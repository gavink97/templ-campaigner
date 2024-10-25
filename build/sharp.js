const sharp = require('sharp');
const fs = require('fs');
const path = require('path');
const process = require('process');

const minifyImages = async (inputDir, outputDir) => {
    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    };

    const images = fs.readdirSync(inputDir)

    for (const image of images) {
        const inputFilePath = path.join(inputDir, image);
        const outputFilePath = path.join(outputDir, image);

        if (fs.existsSync(outputFilePath)) {
            console.log(`Skipping minified image: ${image}`);
            continue;
        }

        if (fs.lstatSync(inputFilePath).isFile()) {
            const ext = path.extname(image).toLowerCase();
            if (ext === '.jpeg' || ext === '.jpg') {
                await minifyJpeg(inputFilePath, outputFilePath)
            } else if (ext === '.png') {
                await minifyPng(inputFilePath, outputFilePath)
            } else {
                console.log(`Skipping unsupported file format: ${image}`)
                continue;
            }
            console.log(`Image minified: ${outputFilePath}`);
        }
    }
};

// look at sharp api to see all options for image minfication
// https://sharp.pixelplumbing.com/api-utility
const minifyJpeg = async (inputFilePath, outputFilePath) => {
    await sharp(inputFilePath)
    .jpeg({ quality: 60 })
    .toFile(outputFilePath);
};

const minifyPng = async (inputFilePath, outputFilePath) => {
    await sharp(inputFilePath)
    .png({ quality: 80, compressionLevel: 9})
    .toFile(outputFilePath)
};

const parseArgs = () => {
    let inputDir = './public/images';
    let outputDir = './bin/images';

    const args = process.argv.slice(2);
    args.forEach((arg, index) => {
        if (arg === '-i' || arg === '--input') {
            inputDir = args[index + 1];
        }
        if (arg === '-o' || arg === '--output') {
            outputDir = args[index + 1];
        }
    });
    return [inputDir, outputDir];
};

const [inputDir, outputDir] = parseArgs();
minifyImages(inputDir, outputDir)
    .then(() => console.log('All images minified successfully.'))
    .catch(err => console.error('Error processing images:', err))
