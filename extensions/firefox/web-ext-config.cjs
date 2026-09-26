module.exports = {
  sourceDir: __dirname,
  artifactsDir: require("node:path").join(__dirname, "dist"),
  ignoreFiles: [
    "node_modules", "node_modules/**", "dist", "dist/**", "tests", "tests/**", "package.json", "package-lock.json",
    "web-ext-config.cjs", "README.md", "*.env", ".env*",
  ],
};
