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

 To run Gator, you'll need the following installed on your system:
 - **Go** 1.16+ (or later)
 - **PostgreSQL** database

 Go programs are statically compiled binaries, so once you've built the application, you can run it without needing the Go toolchain installed.

 ---

 ## 🛠️ Installation & Setup

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
    - Update the configuration file located at `internal/config/config.go` with your database URL and the default user.
      ```go
      type Config struct {
          DBURL           string   Your PostgreSQL connection URL
          CurrentUserName string   Username of the currently logged in user
      }
      ```

 ---

 ## 🚀 Building & Installing the CLI

 For development purposes, you can run the application using:
    ```sh
    go run . <command> [arguments]
    ```

 **Production Build:** Build a statically compiled binary named `gator` that you can run without the Go toolchain.
 To install the CLI globally, run:
    ```sh
    go install github.com/KrishKoria/Gator@latest
    ```
 This will compile and install the `gator` binary into your `$GOPATH/bin` (or `$HOME/go/bin`), which should be in your system PATH.

 ---

 ## 📖 Usage

 After installing, you can run the Gator CLI directly with:
    ```sh
    gator <command> [arguments]
    ```

 ### 🔍 Available Commands

 **User Commands:**
 - `register <username>`: Register a new user.
 - `login <username>`: Log in as an existing user.
 - `users`: List all registered users.
 - `reset`: Reset the database (delete all data).

 **Feed Commands:**
 - `addfeed <name> <url>`: Add a new feed and automatically follow it.
 - `feeds`: List all feeds with user information.
 - `follow <url>`: Follow a feed using its URL.
 - `following`: List all feeds the current user is following.
 - `unfollow <url>`: Unfollow a feed using its URL.
 - `agg <duration>`: Aggregate feed data every specified duration (e.g., "10s", "1m").

 **Post Commands:**
 - `browse [limit]`: Browse posts for the current user, with an optional limit on the number of posts.

 *Note:* `go run .` is intended for development. For production, use the installed `gator` binary.

 ---

 Happy coding and enjoy using Gator! 🐊✨

 If you have any questions or need further assistance, feel free to reach out.
