package repo

import (
	"context"
	"errors"
	"followings/model"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowingRepo struct {
	// Thread-safe instance which maintains a database connection pool
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func New(logger *log.Logger) (*FollowingRepo, error) {
	// Local instance
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &FollowingRepo{
		driver: driver,
		logger: logger,
	}, nil
}

func (mr *FollowingRepo) CheckConnection() {
	ctx := context.Background()
	err := mr.driver.VerifyConnectivity(ctx)
	if err != nil {
		mr.logger.Panic(err)
		return
	}
	// Print Neo4J server address
	mr.logger.Printf(`Neo4J server address: %s`, mr.driver.Target().Host)
}

func (mr *FollowingRepo) CloseDriverConnection(ctx context.Context) {
	mr.driver.Close(ctx)
}

func (mr *FollowingRepo) WritePerson(person *model.User) error {
	// Neo4J Sessions are lightweight so we create one for each transaction (Cassandra sessions are not lightweight!)
	// Sessions are NOT thread safe
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// ExecuteWrite for write transactions (Create/Update/Delete)
	savedPerson, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (p:Person {username: $username}) RETURN p.username + ', from node ' + id(p)",
				map[string]any{"username": person.UserName})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		mr.logger.Println("Error inserting Person:", err)
		return err
	}
	mr.logger.Println(savedPerson.(string))
	return nil
}

// Define a struct to represent the following relationship
type FollowingRelationship struct {
	FollowerUsername string `json:"followerUsername"`
	FollowedUsername string `json:"followedUsername"`
}

// Write a repository method to create the following relationship
func (mr *FollowingRepo) FollowPerson(following *FollowingRelationship) error {
	// Neo4J Sessions are lightweight so we create one for each transaction
	// Sessions are NOT thread safe
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// ExecuteWrite for write transactions (Create/Update/Delete)
	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				"MATCH (follower:Person {username: $followerUsername}) "+
					"MATCH (followed:Person {username: $followedUsername}) "+
					"MERGE (follower)-[:IS_FOLLOWING]->(followed)",
				map[string]any{
					"followerUsername": following.FollowerUsername,
					"followedUsername": following.FollowedUsername,
				})
			return nil, err
		})
	if err != nil {
		mr.logger.Println("Error following person:", err)
		return err
	}
	return nil
}

func (mr *FollowingRepo) GetFollowRecommendations(username string) ([]string, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx) // Provide context.Background() as argument

	// Query to find recommendations
	query := `
		MATCH (u:Person {username: $username})-[:IS_FOLLOWING]->(following)-[:IS_FOLLOWING]->(recommendation:Person)
		WHERE NOT (u)-[:IS_FOLLOWING]->(recommendation)
		RETURN recommendation.username AS username
    `

	// Execute the query
	result, err := session.Run(ctx, query, map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return nil, err
	}

	// Iterate through the results and extract usernames
	var recommendations []string
	for result.Next(ctx) {
		record := result.Record()
		recommendation, ok := record.Get("username")
		if !ok {
			return nil, errors.New("username not found in recommendation")
		}
		recommendations = append(recommendations, recommendation.(string))
		mr.logger.Printf("Found recommendation: %s", recommendation)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}

	return recommendations, nil
}

func (mr *FollowingRepo) GetFollowedUsers(username string) (model.Followed, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (follower:Person {username: $username})-[:IS_FOLLOWING]->(followed:Person)
            RETURN followed.username as username`,
				map[string]interface{}{"username": username})
			if err != nil {
				return nil, err
			}

			var followed model.Followed
			for result.Next(ctx) {
				record := result.Record()
				followedUsername, _ := record.Get("username")
				followed = append(followed, &model.User{UserName: followedUsername.(string)})
			}

			return followed, nil
		})

	if err != nil {
		mr.logger.Println("Error querying followed users:", err)
		return nil, err
	}

	return result.(model.Followed), nil
}

func (mr *FollowingRepo) GetUsersExcept(usernameToExclude string) (model.Users, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (u:Person)
                 WHERE u.username <> $usernameToExclude
                 RETURN u.username AS username`, map[string]interface{}{"usernameToExclude": usernameToExclude})
			if err != nil {
				return nil, err
			}

			var users model.Users
			for result.Next(ctx) {
				record := result.Record()
				username, _ := record.Get("username")
				users = append(users, &model.User{UserName: username.(string)})
			}

			return users, nil
		})

	if err != nil {
		mr.logger.Println("Error querying users:", err)
		return nil, err
	}

	return result.(model.Users), nil
}

func (mr *FollowingRepo) IsFollowing(followerUsername, followedUsername string) (bool, error) {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	query := `
		MATCH (follower:Person {username: $followerUsername})-[:IS_FOLLOWING]->(followed:Person {username: $followedUsername})
		RETURN COUNT(*) > 0 AS isFollowing
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"followerUsername": followerUsername,
		"followedUsername": followedUsername,
	})
	if err != nil {
		return false, err
	}

	if result.Next(ctx) {
		record := result.Record()
		isFollowing, ok := record.Get("isFollowing")
		if !ok {
			return false, errors.New("failed to retrieve isFollowing from result")
		}
		return isFollowing.(bool), nil
	}

	if err := result.Err(); err != nil {
		return false, err
	}

	return false, nil
}
