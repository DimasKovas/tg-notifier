#!/bin/bash
# Token for chat send to. You can get this by /token command in the desirable chat.
CHAT_TOKEN=""
# Default message to send.
MESSAGE="Done"
# Where the bot is hosted. notifier.diko.me is the domain for @simple_notifier_bot
ADDRESS="notifier.diko.me"

printUsage() {
    echo "Usage: notify [options...] [message]"
    echo " -h, --help               Show help"
    echo " -a, --address <address>  Set the host address"
    echo " -t, --token <token>      Set the chat token"
}

while [[ $# -ne 0 ]]; do
    case $1 in
    -h | --help )
        printUsage; exit
        ;;
    -a | --address )
        shift; ADDRESS=$1
        ;;
    -t | --token )
        shift; CHAT_TOKEN=$1
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

if [[ -z "$CHAT_TOKEN" ]]; then
    echo "The chat token should be specified by -t <token>, or by setting the defaults."; exit 1;
fi

curl -d "$MESSAGE" "$ADDRESS/?token=$CHAT_TOKEN" 