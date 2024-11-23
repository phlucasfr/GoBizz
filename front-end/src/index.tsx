/* @refresh reload */
import "./index.css";
import { render } from "solid-js/web";
import App from "./App";

const root = document.getElementById("root");

if (!root) {
  throw new Error("Wrapper div not found");
}

render(() =><App />, root);
