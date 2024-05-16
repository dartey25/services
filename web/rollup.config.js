import css from "rollup-plugin-import-css";
import terser from "@rollup/plugin-terser";
import { nodeResolve } from "@rollup/plugin-node-resolve";
import copy from "rollup-plugin-copy";

export default {
  plugins: [
    nodeResolve(),
    terser(),
    css({
      output: "css/index.min.css",
      minify: true,
    }),
    copy({
      targets: [
        {
          src: "node_modules/@mdoffice/md-ui/fonts/*",
          dest: "../assets/fonts",
        },
      ],
    }),
  ],
  input: { "js/index": "web/src/main.js" },
  external: ["jquery", "htmx.org"],
  output: {
    name: "bundle",
    dir: "../assets",
    entryFileNames: "[name].min.js",
    format: "iife",
    globals: {
      jquery: "$",
      htmx: "htmx",
    },
  },
};
