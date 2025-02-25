 # 🐊 Gator Project

 Gator is a command-line application for managing feeds and follows. With Gator, you can easily register, log in, add feeds, follow/unfollow feeds, and view the feeds you care about—all powered by a PostgreSQL database!
 
 ---
 
 ## 🚀 Features
 
 - **User Management**: Register and log in with ease.
 - **Feed Management**: Add new feeds effortlessly.
 - **Following System**: Follow or unfollow feeds with a simple command.
 - **Data Viewing**: List all feeds with associated user information.
 - **Personalized Feed**: Quickly view feeds you're following.
 
 ---
 
 ## 📋 Prerequisites
 
 Before getting started, make sure you have the following:
 
 - **Go** 1.16+ installed
 - **PostgreSQL** database installed and running
 - Required Go packages:
   - [`github.com/lib/pq`](https:github.com/lib/pq) for PostgreSQL driver
   - [`github.com/google/uuid`](https:github.com/google/uuid) for UUID generation
 
 ---
 
 ## 🛠️ Installation
 
 Follow these steps to set up Gator on your machine:
 
 1. **Clone the Repository:**
    ```sh
    git clone https:github.com/KrishKoria/Gator.git
    cd Gator
    ```
 
 2. **Install Dependencies:**
    ```sh
    go mod tidy
    ```
 
 3. **Configure Your Database:**
    - Set up your PostgreSQL database.
    - Update the configuration file (`internal/config/config.go`) with your database URL.
 
 ---
 
 ## ⚙️ Configuration
 
 Gator uses a configuration file to manage its settings. Update the file at `internal/config/config.go` with the following structure:
 
 ```go
 type Config struct {
     DBURL           string  Your PostgreSQL connection URL
     CurrentUserName string  Username of the currently logged in user
 }
 ```
 
 Ensure that your `DBURL` is correctly set so that Gator can connect to your PostgreSQL instance.
 
 ---
 
 ## 📖 Usage
 
 Run the application by navigating to the project directory and executing:
 
 ```sh
 go run . <command> [arguments]
 ```
 
 ### 🔍 Available Commands
 
 - **User Commands:**
   - `register <username>`: Register a new user.
   - `login <username>`: Log in as an existing user.
   - `users`: List all registered users.
   - `reset`: Reset the database (delete all users).
 
 - **Feed Commands:**
   - `addfeed <name> <url>`: Add a new feed and automatically follow it.
   - `feeds`: List all feeds with user information.
   - `follow <url>`: Follow a feed using its URL.
   - `following`: List all feeds the current user is following.
   - `unfollow <url>`: Unfollow a feed using its URL.
   - `agg`: Aggregate feed data.
 
 ---
 
 ## 🗄️ SQL Queries
 
 The SQL queries used for managing feeds are located in the `sql/queries/feeds.sql` file. They include:
 
 - **CreateFeed**: Insert a new feed.
 - **GetFeedsWithUsers**: Retrieve all feeds along with user details.
 - **GetFeedByURL**: Fetch a feed by its URL.
 - **CreateFeedFollow**: Add a new feed follow record and return related information.
 - **GetFeedFollowsForUser**: Get all feed follow records for a user.
 - **DeleteFeedFollowByUserAndFeedURL**: Remove a feed follow record using the user and feed URL.
 
 ---
 
 ## 🛡️ Middleware
 
 To ensure that only authenticated users perform certain actions, Gator uses middleware. For example, the `middlewareLoggedIn` function verifies the current user before processing commands:
 
 ```go
 func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
     return func(s *state, cmd command) error {
         currentUser := s.Config.CurrentUserName
         user, err := s.DBQueries.GetUser(context.Background(), currentUser)
         if err != nil {
             return fmt.Errorf("error getting current user: %v", err)
         }
         return handler(s, cmd, user)
     }
 }
 ```
 
 ---
 
 ## 📝 Handlers
 
 The command handlers in `handlers.go` define the behavior for each command:
 
 - **User Handlers:**
   - `handlerLogin`: Processes user login.
   - `handlerRegister`: Handles user registration.
   - `handlerUsers`: Displays all registered users.
   - `handlerReset`: Resets the database.
 
 - **Feed Handlers:**
   - `handlerAddFeed`: Adds a new feed and follows it automatically.
   - `handlerFeeds`: Lists all feeds with user details.
   - `handlerFollow`: Follows a feed using its URL.
   - `handlerFollowing`: Lists feeds followed by the current user.
   - `handlerUnfollow`: Unfollows a feed using its URL.
   - `handlerAgg`: Aggregates feed data (example command).
 
 ---
 
 ## 📜 License
 
 This project is licensed under the **MIT License**. Feel free to use, modify, and distribute it as needed.
 
 ---
 
 Happy coding and enjoy using Gator! 🐊✨
 
 Feel free to reach out if you have any questions or need further enhancements.
