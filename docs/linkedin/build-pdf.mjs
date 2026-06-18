// Builds a single 6-page square PDF source from slides.html.
// Adds an @page rule (1080x1080) so each .slide is one square page.
import { readFileSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = dirname(fileURLToPath(import.meta.url));
const html = readFileSync(join(here, 'slides.html'), 'utf8');

const out = html.replace(
  '</style>',
  '@page { size: 1080px 1080px; margin: 0; }\n</style>'
);

writeFileSync(join(here, 'pdf', '_carousel.html'), out);
console.log('Wrote pdf/_carousel.html');
