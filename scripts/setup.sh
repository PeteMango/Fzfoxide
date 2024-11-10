#!/bin/bash

set -e

# Define variables
ZSHRC="$HOME/.zshrc"
CD_FUNCTION='
cd() {
    if [ "$#" -eq 0 ]; then
        builtin cd ~
        fzfoxide --run ~ > /dev/null
    else
        if builtin cd "$1" 2>/dev/null; then
            fzfoxide --run "$(pwd)" > /dev/null
        else
            dir=$(fzfoxide --run "$1")
            if [ -n "$dir" ] && [ -d "$dir" ]; then
                builtin cd "$dir"
            else
                echo "No matching directory found"
            fi
        fi
    fi
}
'

add_cd_function() {
    echo "$CD_FUNCTION" >> "$ZSHRC"
    echo "Custom cd function added to $ZSHRC."
}

install_fzfoxide() {
    echo "Compiling fzfoxide Go program..."

    mkdir -p bin

    go build -o bin/fzfoxide cmd/fzfoxide/main.go

    if [ ! -f "bin/fzfoxide" ]; then
        echo "Error: fzfoxide binary was not created. Please check your Go code for errors."
        exit 1
    fi

    echo "fzfoxide compiled successfully."

    echo "Installing fzfoxide to /usr/local/bin/..."

    if [ ! -d "/usr/local/bin" ]; then
        echo "Error: /usr/local/bin directory does not exist."
        exit 1
    fi

    sudo mv bin/fzfoxide /usr/local/bin/

    sudo chmod +x /usr/local/bin/fzfoxide

    echo "fzfoxide installed to /usr/local/bin/ successfully."
}

main() {
    echo "Starting setup..."

    add_cd_function

    install_fzfoxide

    echo "Setup completed successfully!"
    echo "Please restart your terminal or open a new one to start using the updated cd function."
    echo "Alternatively, you can source your .zshrc manually by running 'source ~/.zshrc' within a Zsh session."
}

main
