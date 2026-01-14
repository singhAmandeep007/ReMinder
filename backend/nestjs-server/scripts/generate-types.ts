/**
 * Generate TypeScript types from OpenAPI specification
 *
 * This script generates TypeScript types from the OpenAPI spec
 * located at /openapi/api.yaml. The generated types are placed
 * in src/generated/api.types.ts.
 *
 * Usage: npx ts-node scripts/generate-types.ts
 * Or via npm: npm run generate:types
 */

import { execSync } from 'child_process';
import * as path from 'path';
import * as fs from 'fs';

const OPENAPI_SPEC = path.resolve(__dirname, '../../../openapi/api.yaml');
const OUTPUT_DIR = path.resolve(__dirname, '../src/generated');
const OUTPUT_FILE = path.join(OUTPUT_DIR, 'api.types.ts');

// Ensure output directory exists
if (!fs.existsSync(OUTPUT_DIR)) {
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

// Check if OpenAPI spec exists
if (!fs.existsSync(OPENAPI_SPEC)) {
  console.error(`❌ OpenAPI spec not found at: ${OPENAPI_SPEC}`);
  process.exit(1);
}

console.log('🔧 Generating TypeScript types from OpenAPI spec...');
console.log(`   Source: ${OPENAPI_SPEC}`);
console.log(`   Output: ${OUTPUT_FILE}`);

try {
  execSync(`npx openapi-typescript ${OPENAPI_SPEC} -o ${OUTPUT_FILE}`, {
    stdio: 'inherit',
  });
  console.log('✅ Types generated successfully!');
} catch (error) {
  console.error('❌ Failed to generate types:', error);
  process.exit(1);
}
