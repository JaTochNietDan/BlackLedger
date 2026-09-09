import {defineConfig} from 'vite';
import {resolve} from 'path';

// Two pages out of one build.
//
// The game is index.html. playtest.html is the workshop: the same components,
// driven by hand instead of by a campaign, so the visual moments can be looked
// at without playing far enough to make one happen. It is served by its own
// binary on its own port and is not part of the game — that is the whole point
// of it being a separate entry rather than a route.
export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        playtest: resolve(__dirname, 'playtest.html'),
      },
    },
  },
});
