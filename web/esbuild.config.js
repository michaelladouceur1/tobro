const esbuild = require("esbuild");
const svgrPlugin = require("esbuild-plugin-svgr");

async function watch() {
  let ctx = await esbuild.context({
    entryPoints: ["./src/index.tsx"],
    bundle: true,
    minify: true,
    format: "cjs",
    sourcemap: false,
    outfile: "./build/output.js",
    plugins: [svgrPlugin()],
    loader: { ".ts": "tsx", ".js": "jsx", ".svg": "file" },
  });
  await ctx.watch();
  console.log('Watching...');
}
watch();