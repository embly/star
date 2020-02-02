import Terminal from "xterm/src/xterm.ts";
import * as fullscreen from "xterm/src/addons/fullscreen/fullscreen.ts";
import * as fit from "xterm/src/addons/fit/fit.ts";

import "xterm/dist/xterm.css";
import "xterm/dist/addons/fullscreen/fullscreen.css";

// imports "Go"
import "./wasm_exec.js";

Terminal.applyAddon(fullscreen);
Terminal.applyAddon(fit);

// Polyfill for WebAssembly on Safari
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

const term = new Terminal();
term.open(document.getElementById("terminal"));
term.toggleFullScreen();
term.fit();

term.writeln("loading...");
term.writeln("(downloading and parsing wasm binary, this could take a while)");
window.onresize = () => {
  term.fit();
};

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
  result => {
    let mod = result.module;
    let inst = result.instance;
    go.run(inst);
  }
);

let prompt_prefix = ">>> ";

window.set_prompt = val => {
  prompt_prefix = val;
};

const up = 38;
const down = 40;

function prompt(term) {
  term.write("\r\n" + prompt_prefix);
}

window.term_write = val => {
  term.write("\r\n");
  term.write(val);
};

term.attachCustomKeyEventHandler(e => {
  let out = e.key !== "v" && !e.ctrlKey;
  return out;
});
// called when go is done initializing
window.go_ready = () => {
  term.clear();
  term.writeln("Welcome to the star repl.");
  term.writeln(
    'Try writing some python or import a Go lib with require("net/http")'
  );
  term.prompt = () => {
    term.write("\r\n" + prompt_prefix);
  };

  prompt(term);
  var buffer = [];
  term.onKey(e => {
    const printable =
      !e.domEvent.altKey &&
      !e.domEvent.altGraphKey &&
      !e.domEvent.ctrlKey &&
      !e.domEvent.metaKey;

    if (e.domEvent.keyCode === up || e.domEvent.keyCode == down) {
    } else if (e.domEvent.keyCode === 13) {
      newLine(buffer.join(""));
      buffer = [];
      prompt(term);
    } else if (e.domEvent.keyCode === 8) {
      // Do not delete the prompt
      if (term._core.buffer.x > 4) {
        buffer.pop();
        term.write("\b \b");
      }
    } else if (printable) {
      buffer.push(e.key);
      term.write(e.key);
    }
  });
};
