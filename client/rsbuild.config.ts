import { defineConfig } from '@rsbuild/core';
import { pluginSvelte } from '@rsbuild/plugin-svelte';
import { pluginSass } from '@rsbuild/plugin-sass';
// Docs: https://rsbuild.rs/config/
export default defineConfig({
  plugins: [pluginSvelte(),pluginSass()],
});
