#! /bin/sh

prettier -w ./script.js && prettier -w ./index.html && prettier -w style.css && npx eslint ./script.js
