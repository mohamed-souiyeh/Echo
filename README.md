# Echo

Simple terminal chat application served via SSH - which means u only need a terminal and an ssh client to start using it!


## ⚠️ Status: Work In Progress ⚠️

**Echo is currently under active development.** Core features like chat rooms and message persistence are not yet implemented. Authentication is the current focus. Expect bugs and incomplete functionality.

## Features

*   **SSH Access:** Connect using any standard SSH client.
*   **Terminal UI:** Built using the lovely Bubble Tea framework.
*   **(In Progress):** User Authentication (Signup/Login).
*   **(Planned):**
    *   Direct Messaging (DMs) / Rooms
    *   Persistent Chat History (using SQLite)
    *   User Search
    *   Real-time message broadcasting

## Tech Stack

*   **Language:** Go
*   **SSH Server:** `charmbracelet/wish`
*   **Terminal UI (TUI):** `charmbracelet/bubbletea`
*   **TUI Styling:** `charmbracelet/lipgloss`
*   **Database:** SQLite (Embedded)
*   **DB Migrations:** `golang-migrate/migrate`
*   **DB Query Generation:** `sqlc`

## Getting Started (Running Locally)

These instructions are for developers or testers wanting to run the server locally.

1.  **Prerequisites:**
    *   Go (version 1.24.1 or later)
    *   An SSH client (`ssh`)

2.  **Clone the repository:**
    ```bash
    git clone https://github.com/mohamed-souiyeh/Echo
    cd Echo
    ```

3.  **Set up the Database:**
    *   The application uses an embedded SQLite database (`echo.db` by default for now).
    *   The database migration is run automatically ont he server startup.

4.  **Generate SQLC Code:**
    *   (If you modify `./db/queries/*.sql` files) Regenerate the Go database code:
        ```bash
        sqlc generate
        ```

5.  **Build the application:**
    ```bash
    go build -o echo .
    ```

6.  **Run the server:**
    ```bash
    ./echo
    # It should start listening for SSH connections on port 4242 by default.
    ```

## Usage (Connecting)

Once the server is running locally:

1.  Connect using your SSH client:
    ```bash
    ssh localhost -p 4242
    ```
2.  You should be presented with the authentication screen (Login/Signup). Follow the on-screen prompts.


