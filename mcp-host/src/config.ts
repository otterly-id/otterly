import { readFileSync } from 'fs';
import { load } from 'js-yaml';

interface Config {
  API_KEY: string;
  GEMINI_API_KEY: string;
}

const loadConfig = (): Config => {
  try {
    const configPath = './config/env.yaml';
    const fileContents = readFileSync(configPath, 'utf8');
    const config = load(fileContents) as Config;
    return config;
  } catch (e) {
    console.error('Error loading configuration:', e);
    process.exit(1);
  }
};

export const config = loadConfig(); 