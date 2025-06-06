#!/usr/bin/bash

# Default values for the flags
FLAG_CCAPI="none"
FLAG_LABEL="1.0"

# You can change this if you want to avoid using the --name flag
FLAG_NAME="clausia-cc"

# Process command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --ccapi | -c)
            FLAG_CCAPI="go"
            shift 1
            ;;
        --label | -l)
            if [[ $# -gt 1 ]]; then
                FLAG_LABEL=$2
                shift 2
            else
                echo "Error: --label flag requires a value."
                exit 1
            fi
            ;;
        --name | -n)
            if [[ $# -gt 1 ]]; then
                FLAG_NAME=$2
                shift 2
            else
                echo "Error: --name flag requires a value."
                exit 1
            fi
            ;;
        --help | -h)
            echo "Usage: ./generateTar.sh [--ccapi] [--label <label>] [--name <name>]"
            echo "  --ccapi, -c: Include rest-server in the tar file. Default is no ccapi."
            echo "  --help , -h: Show this help message."
            echo "  --label, -l: Label to be used for the chaincode package. Default is 1.0."
            echo "  --name , -n: Name of the chaincode package. Default is ${FLAG_NAME}."
            exit 0
            ;;
        *)
            # Ignore unrecognized arguments
            shift
            ;;
    esac
done

# Remove previous tar file
rm -f ${FLAG_NAME}.tar.gz

# Make sure go mod is up to date
cd chaincode && go mod vendor && cd ..

# Pack chaincode
export FABRIC_CFG_PATH=fabric/config
peer lifecycle chaincode package chaincode.tar.gz --path chaincode --lang golang --label ${FLAG_NAME}_${FLAG_LABEL}

case $FLAG_CCAPI in
    go)
        # Create temporary directory for rest-server
        mkdir -p tmp/rest-server
        cp -r ccapi/* tmp/rest-server

        # Compress file with rest-server (GoFabric will use the one provided)
        tar -c --exclude=vendor -zf ${FLAG_NAME}.tar.gz chaincode.tar.gz -C tmp rest-server
        ;;
    none)
        # Compress file without rest-server (GoFabric will use the standard CC API)
        tar -czf ${FLAG_NAME}.tar.gz chaincode.tar.gz
esac

# Remove chaincode.tar.gz
rm -f chaincode.tar.gz

# Remove temporary directory
rm -rf tmp