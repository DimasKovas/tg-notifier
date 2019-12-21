#!/bin/bash
# Where the bot is hosted. notifier.diko.me if the domain for @simple_notifier_bot
ADDRESS="notifier.diko.me"
# Which message will be sent to you.
MESSAGE="Done"
# To which chat to send.
USER_ID=""

printUsage() {
    echo "Usage: notify [options...] [message]"
    echo " -h, --help               Show help"
    echo " -a, --address <address>  Host address"
    echo " -u, --userid <userid>    Chat ID"
}

while [[ $# -ne 0 ]]; do
    case $1 in
    -h | --help )
        printUsage; exit
        ;;
    -a | --address )
        shift; ADDRESS=$1
        ;;
    -u | --userid )
        shift; USER_ID=$1
        ;;
    -*)
        echo "Unexpected option '$1'"; printUsage; exit 1
        ;;
    *)
        MESSAGE=$1
        ;;
    esac;
    shift;
done

if [[ -z "$USER_ID" ]]; then
    echo "User ID should be specified by -u <userid>, or be setting defaults."; exit 1;
fi

curl -H "UserID: $USER_ID" -d "$MESSAGE" "$ADDRESS" 