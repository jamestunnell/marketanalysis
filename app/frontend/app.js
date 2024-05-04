import van from "vanjs-core"
import {Securities} from "./securities.js"

van.add(document.body, Securities())

console.log("backend URL: %s", process.env.BACKEND_URL);