# Overview
This project is structured to facilitate modular development and oepration.
The fundimental part of this project is the root project itself should only provide the basic input and output handling for the bot.
THis means that the root project should manage the connection to the IRC server and reading/writing to the postgresql database, but the actual logic of the bot should be implemented in separate modules.
This allows for a clean separation of concerns and makes it easier to maintain and extend the bot's functionality.

The configuration for the bots should be stored in the postgresql database. This allows for dynamic configuration changes without needing to restart the bot, and also provides a centralized location for managing the bot's settings. The configuration can include things like the bot's nickname, the channels it should join, and any other settings that are relevant to the bot's operation. By storing this information in the database, we can easily update the bot's configuration as needed without having to modify the code or restart the bot.

The modules should be located in a subfolder like modules/PMProxy

## Bot config
The bot configuration should include:
- Connection options for the IRC server (e.g., server address, port, SSL settings)
- Bot nickname and username
- List of channels to join
- The adminsitration (botteam) channel
- the log channel
- Any other relevant settings for the bot's operation this should be able to be dynamically set by individual modules as well, for example the PMProxy module should be able to set the channel it uses for private message proxying without needing to restart the bot.

## Process management
The bot should be designed to run as a long-running process, and should be able to handle unexpected errors gracefully. This means that if the bot encounters an error, it should log the error and attempt to recover without crashing. Additionally, the bot should be able to restart itself if it crashes, either through a built-in mechanism or by using an external process manager like systemd or supervisord. This ensures that the bot remains operational even in the face of unexpected issues.

Each process should only run 1 bot, multiple instances of program will be used to run multiple bots. 
The ID of the bot being ran will be given as an argument to the program, and the program will read the configuration for that bot from the database.
