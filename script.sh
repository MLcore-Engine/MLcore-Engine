#!/bin/bash

run_script() {
    cd web/ || { echo "Failed to cd into web/"; exit 1;}
    echo "Running npm build..."
    npm run build || { echo "npm run build failed"; exit 1;}
    cd .. || { echo "Failed to cd .."; exit 1;}
    echo "Running Go program..."
    go run main.go || { echo "go run main.go failed"; exit 1;}
}

install_dependencies() {
    package_name=$1
    if [ -z "$package_name" ]; then
        echo "Package name is required"
        exit 1
    fi
    echo "cd web/"
    cd web/ || { echo "Failed to cd into web/"; exit 1;}
    echo "Installing $package_name..."
    npm install $package_name || { echo "npm install  $package_name failed"; exit 1;}
    cd .. || { echo "Failed to cd .."; exit 1;}
}

case $1 in
    run)
        run_script
        ;;
    install)
        if [ -z "$2" ]; then
            echo "Usage: $0 install <package_name>"
            exit 1;
        fi
        install_dependencies "$2"
        ;;
    *)
        echo "Usage: $0 {run|install <package_name>}"
        exit 1
        ;;
esac

