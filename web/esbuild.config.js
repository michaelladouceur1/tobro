// const res = require('esbuild').buildSync({
//     entryPoints: ['src/index.tsx'],
//     bundle: true,
//     minify: true,
//     format: 'cjs',
//     sourcemap: true,
//     outfile: 'build/output.js',
//     // external: ['react', 'react-dom'], 
//   })

const esbuild = require("esbuild");
async function watch() {
  let ctx = await esbuild.context({
    entryPoints: ["./src/index.tsx"],
    bundle: true,
    minify: true,
    format: "cjs",
    sourcemap: true,
    outfile: "./build/output.js",
    loader: { ".ts": "ts" },
  });
  await ctx.watch();
  console.log('Watching...');
}
watch();