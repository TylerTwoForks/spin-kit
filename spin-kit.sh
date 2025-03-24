#!/bin/bash

## customize these values #############################
sandbox_list=("update_me1" "update_me2" "update_me3")
prod_alias="Prod"
#######################################################

index=${#sandbox_list[@]}
trap 'menu' ERR

enterChoice() {
    unset $choice
    unset $sandbox_name
    read -p "Enter your choice: " choice
}

#calculates minimum width needed and adds a bit of padding
calculateWidth() {
    max_length=0
    for i in "${!sandbox_list[@]}"; do
        length=${#sandbox_list[$i]}
        if [ $length -gt $max_length ]; then
            max_length=$length
        fi
    done
    width=$((max_length + 22))
}

showMenu() {
    unset $choice
    unset $sandbox_name
    calculateWidth
    echo "┌$(printf '─%.0s' $(seq 1 $width))┐"
    echo "│$(printf ' %.0s' $(seq 1 $((width / 2 - 5))))MAIN MENU$(printf ' %.0s' $(seq 1 $((width / 2 - 4))))│"
    echo "├$(printf '─%.0s' $(seq 1 $width))┤"
    echo "│ Refresh $(printf ' %.0s' $(seq 1 $((width - 9))))│"
    for i in "${!sandbox_list[@]}"; do
        printf "│  %d: %s %$(($width - ${#sandbox_list[$i]} - 6))s│\n" "$i" "${sandbox_list[$i]}" ""
    done
    echo "│ $(printf ' %.0s' $(seq 1 $((width - 1))))│"
    echo "│ Utilities $(printf ' %.0s' $(seq 1 $((width - 11))))│"
    echo -e "│  ls:\tList Org Connections $(printf ' %.0s' $(seq 1 $((width - 28))))│"
    echo -e "│  ref:\tRefresh Custom Sandbox $(printf ' %.0s' $(seq 1 $((width - 30))))│"
    echo -e "│  rc:\tReconnect to Sandbox $(printf ' %.0s' $(seq 1 $((width - 28))))│"
    echo -e "│  m:\tShow menu $(printf ' %.0s' $(seq 1 $((width - 17))))│"
    echo -e "│  x:\tExit $(printf ' %.0s' $(seq 1 $((width - 12))))│"
    echo -e "└$(printf '─%.0s' $(seq 1 $width))┘"

    enterChoice
}

isValidChoice() {
    if [[ "$sandbox_name" == "" ]]; then
        unset $sandbox_name
        return 1
    elif [[ "$sandbox_name" == "exit" || "$sandbox_name" == "x" ]]; then
        unset $sandbox_name
        echo "returning to menu..."
        myApp
        return 0
    else
        return 0
    fi
}

myApp() {
    unset $choice
    unset $sandbox_name
    showMenu
    while true; do
        # dynamic choices
        if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 0 ] && [ "$choice" -lt "$index" ]; then
            echo "Refreshing ${sandbox_list[$choice]}"
            sf org refresh sandbox --name ${sandbox_list[$choice]} --target-org $prod_alias
            # continue
        fi

        # static choices
        case $choice in
        "")
            myApp
            ;;
        ls)
            echo "Org Connections:"
            sf org list #sf command
            ;;
        x | exit)
            exit
            ;;
        m)
            myApp
            ;;
        ref)
            echo
            read -p "Enter Sandbox Name to Refresh (or "x" to return to menu): " sandbox_name
            if ! isValidChoice; then
                echo "Invalid choice, please try again."
            else
                echo "^c to return to menu or continue the login process in the browser"
                sf org refresh sandbox --name $sandbox_name --target-org $prod_alias #sf command
                
            fi
            ;;
        rc)
            echo
            read -p "Enter Sandbox Name to Reconnect (or "x" to return to menu): " sandbox_name
            if ! isValidChoice; then
                echo "Invalid choice, please try again."
            else
                echo "^c to return to menu or continue the login process in the browser"
                sf org web login --instance-url https://test.salesforce.com --alias $sandbox_name #sf command
                
            fi

            ;;
        *)
            echo "Invalid choice, please try again."
            ;;
        esac
        myApp

    done
}

myApp
