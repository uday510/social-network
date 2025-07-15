package db

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/uday510/go-crud-app/internal/store"
	"log"
	"time"
)

func Seed(store store.Storage) error {
	ctx := context.Background()

	log.Println("Generating users...")
	users := generateUsers(1000)
	log.Printf("Generated %d users", len(users))

	log.Println("Inserting users into DB...")
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("Error creating user: %s\n", err)
			return err
		}
	}
	log.Println("Inserted users")

	log.Println("Generating posts...")
	posts := generatePosts(2000, users)
	log.Printf("Generated %d posts", len(posts))

	log.Println("Inserting posts into DB...")
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("Error creating post: %s\n", err)
			return err
		}
	}
	log.Println("Inserted posts")

	log.Println("Generating comments...")
	comments := generateComments(10000, posts, users)
	log.Printf("Generated %d comments", len(comments))

	log.Println("Inserting comments into DB...")
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("Error creating comment: %s\n", err)
			return err
		}
	}
	log.Println("Inserted comments")

	log.Println("Seeding complete!")
	return nil
}

func generateUsers(num int) []*store.User {
	gofakeit.Seed(time.Now().UnixNano())

	users := make([]*store.User, 0, num)
	usedUsernames := make(map[string]bool)
	usedEmails := make(map[string]bool)

	for len(users) < num {
		username := gofakeit.Username()
		email := gofakeit.Email()

		if usedUsernames[username] || usedEmails[email] {
			continue
		}

		usedUsernames[username] = true
		usedEmails[email] = true

		user := &store.User{
			Username: username,
			Email:    email,
			Password: gofakeit.Password(true, true, true, true, false, 12),
		}
		users = append(users, user)
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[gofakeit.Number(0, len(users)-1)]
		posts[i] = &store.Post{
			Title:    gofakeit.Sentence(6),
			Content:  gofakeit.Paragraph(1, 3, 12, " "),
			UserID:   user.ID,
			Tags:     randomWords(gofakeit.Number(1, 5)),
			Version:  1,
			Comments: []store.Comment{},
		}
	}
	return posts
}

func generateComments(num int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, 0, num)
	for i := 0; i < num; i++ {
		post := posts[gofakeit.Number(0, len(posts)-1)]
		user := users[gofakeit.Number(0, len(users)-1)]
		comment := &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: gofakeit.Sentence(10),
		}
		comments = append(comments, comment)
		post.Comments = append(post.Comments, *comment)
	}
	return comments
}

func randomWords(n int) []string {
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = gofakeit.Word()
	}
	return words
}
