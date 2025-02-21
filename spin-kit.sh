#!/bin/bash

## customize these values #############################
sandbox_list=("zincTe2Dev" "update_me2" "update_me3")
prod_alias="Prod"
#######################################################

index=${#sandbox_list[@]}
trap 'menu' ERR

enterChoice(){
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
    width=$((max_length + 30))
}

menu() {
    calculateWidth
    echo "┌$(printf '─%.0s' $(seq 1 $width))┐"
    echo "│$(printf ' %.0s' $(seq 1 $((width / 2 - 5))))MAIN MENU$(printf ' %.0s' $(seq 1 $((width / 2 - 4))))│"
    echo "├$(printf '─%.0s' $(seq 1 $width))┤"
    for i in "${!sandbox_list[@]}"; do
        printf "│ %d: Refresh %s %$(($width - ${#sandbox_list[$i]} - 13))s│\n" "$i" "${sandbox_list[$i]}" ""
        echo "├$(printf '─%.0s' $(seq 1 $width))┤"
    done
    echo "│ l | ls: List Org Connections $(printf ' %.0s' $(seq 1 $((width - 30))))│"
    echo "├$(printf '─%.0s' $(seq 1 $width))┤"
    echo "│ c | custom: Refresh Custom Sandbox $(printf ' %.0s' $(seq 1 $((width - 36))))│"
    echo "├$(printf '─%.0s' $(seq 1 $width))┤"
    echo "│ m | menu: Show menu $(printf ' %.0s' $(seq 1 $((width - 21))))│"
    echo "├$(printf '─%.0s' $(seq 1 $width))┤"
    echo "│ x | exit: Exit $(printf ' %.0s' $(seq 1 $((width - 16))))│"
    echo "└$(printf '─%.0s' $(seq 1 $width))┘"

    enterChoice
}

orgStatus(){
    echo "Org Connections:";
    sf org list;
    enterChoice
}



menu
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
            menu
            ;;
        ls|list|l) 
            orgStatus
            ;;
        x)
            exit
            ;;
        m|menu)
            menu
            ;;
        c|custom)
            echo
            read -p "Enter Sandbox Name to Refresh: " sandbox_name
            sf org refresh sandbox --name $sandbox_name --target-org $prod_alias
            enterChoice
            ;;
        *)
            echo "Invalid choice, please try again."
            enterChoice
            ;;
        
    esac;

done