# 🎵 Code & Play: Automate Audio Notifications for Your Commands 🎵

Tired of staring at your terminal, waiting for commands to finish? Let’s make your workflow more fun and efficient by playing an audio notification when your tasks are done! Here’s how to use the **Ring API** to turn your terminal into a DJ booth.

---

## 🛠️ What You’ll Need
1. **Ring API** running (local or remote)
2. A terminal (bash, zsh, etc.)
3. A cool audio file (e.g., `success.wav`)
4. A sense of adventure 🚀

---

## 🎧 Step 1: Set Up the Ring API
First, start the Ring API service:

```bash
# Clone the repo
git clone https://github.com/pluveto/ring-api.git
cd ring-api

# Build and run
AUDIO_PATH=/path/to/your/audio.wav go run cmd/ring-api/main.go
```

Or, use Docker for instant setup:

```bash
docker run -p 8080:8080 -e AUDIO_PATH=/audio/success.wav pluveto/ring-api
```

---

## 🎯 Step 2: Create a Command Wrapper
Write a bash function to wrap your commands and play audio when they finish.

Add this to your `.bashrc` or `.zshrc`:

```bash
function notify-me() {
    # Run the command
    "$@"
    
    # Check if the command succeeded
    if [ $? -eq 0 ]; then
        # Play the success sound
        curl -X GET http://localhost:8080/api/ring
    else
        echo "Command failed. No music for you! 🥲"
    fi
}
```

Reload your shell config:

```bash
source ~/.bashrc  # or ~/.zshrc
```

---

## 🚀 Step 3: Use It!
Now, wrap any command with `notify-me` and enjoy the music when it’s done!

```bash
# Example: Build your Go project
notify-me go build ./...

# Example: Run tests
notify-me go test ./...

# Example: Long-running task
notify-me sleep 10
```

When the command finishes, the Ring API will play your audio file! 🎶

---

## 🎨 Bonus: Customize Your Experience
1. **Change the Audio**: Replace `success.wav` with your favorite sound (e.g., a victory fanfare 🎺 or a meme sound 🐸).
2. **Remote API**: Deploy the Ring API to a server and update the `curl` command to point to your remote URL.
3. **Error Sounds**: Add a different audio file for failed commands.

---

## 🌟 Example Workflow
```bash
# Start your day with a build
notify-me make build

# Run your tests
notify-me make test

# Deploy your app
notify-me make deploy
```

Now you can step away from your terminal, grab a coffee ☕, and let the music tell you when your work is done!

---

## 🎉 Why Wait? Start Coding & Playing Today!
With the Ring API, you’ll never miss a command completion again. Turn your terminal into a productivity DJ and make coding more fun! 🎧✨

---

**Pro Tip**: Use a short, distinctive sound so you know exactly when your task is done without even looking at the terminal. 🎶
