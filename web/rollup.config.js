import css from "rollup-plugin-import-css";
import terser from "@rollup/plugin-terser";
import resolve from "@rollup/plugin-node-resolve";
import copy from "rollup-plugin-copy";
import commonjs from "@rollup/plugin-commonjs";

export default {
  plugins: [
    resolve(),
    commonjs(),
    terser(),
    css({
      output: "css/index.min.css",
      minify: true,
    }),
    copy({
      targets: [
        {
          src: "web/node_modules/@mdoffice/md-ui/fonts/*",
          dest: "assets/fonts",
        },
      ],
    }),
  ],
  input: { "js/index": "web/src/js/main.js" },
  output: {
    name: "bundle",
    dir: "assets",
    entryFileNames: "[name].min.js",
    format: "iife",
  },
};
