import { resolve, join } from "path";
import { readdirSync, statSync } from "fs";

import { defineConfig } from "vite";

const root = resolve(__dirname, "src");
const outDir = resolve(__dirname, "dist");
const inputOpts: { [p: string]: string } = {
    main: resolve(root, "index.html"),
};

// Jetbrains suppression syntax for unused export default
// noinspection JSUnusedGlobalSymbols
export default defineConfig({
    root,
    build: {
        outDir,
        emptyOutDir: true,
        rollupOptions: {
            input: inputOpts,
        },
    },
});

function getSrcFiles(dirPath: string, arrayOfFiles: string[] = []) {
    const files = readdirSync(dirPath);

    for (const file of files) {
        if (file === "public") {
            continue;
        }
        if (statSync(dirPath + "/" + file).isDirectory()) {
            arrayOfFiles = getSrcFiles(
                dirPath + "/" + file,
                arrayOfFiles,
            );
            continue;
        }
        if (file !== "index.html") {
            continue;
        }
        const srcFilePath = join(dirPath, "/", file);
        arrayOfFiles.push(srcFilePath);
    }

    return arrayOfFiles;
}

for (const filePath of getSrcFiles(root)) {
    const dirs = filePath.split("/");
    const directParent = dirs[dirs.length - 2];
    if (directParent === "src") {
        continue;
    }
    inputOpts[directParent] = filePath;
}
