/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

// See https://docusaurus.io/docs/site-config for all the possible
// site configuration options.

const siteConfig = {
  title: 'Graphql Terraform Provider',

  url: 'https://sullivtr.github.io',
  baseUrl: '/terraform-provider-graphql/',
  projectName: 'terraform-provider-graphql',
  organizationName: 'sullivtr',

  editUrl: 'https://github.com/sullivtr/terraform-provider-graphql/edit/master/docusaurus/docs/',

  // For no header links in the top nav bar -> headerLinks: [],
  headerLinks: [
    { href: "https://github.com/sullivtr/terraform-provider-graphql/releases", label: "Releases" },
    { href: "https://github.com/sullivtr/terraform-provider-graphql", label: "GitHub" },
  ],

  /* path to images for header/footer */
  favicon: 'img/favicon.png',

  /* Colors for website */
  colors: {
    primaryColor: '#171e26',
    secondaryColor: '#E10098',
  },

  usePrism: ['yaml', 'js', 'bash', 'sh', 'hcl'],

  // This copyright info is used in /core/Footer.js and blog RSS/Atom feeds.
  copyright: `Copyright © ${new Date().getFullYear()} Tyler R. Sullivan`,

  highlight: {
    // Highlight.js theme to use for syntax highlighting in code blocks.
    theme: 'default',
  },

  // Add custom scripts here that would be placed in <script> tags.
  scripts: ['https://buttons.github.io/buttons.js'],

  stylesheets: [
    "https://fonts.googleapis.com/css?family=Open Sans:400,400i,700,700i,900,900i&display=swap"
  ],

  // On page navigation for the current documentation page.
  onPageNav: 'separate',

  // No .html extensions for paths.
  cleanUrl: true,
  docsSideNavCollapsible: true,
};

module.exports = siteConfig;
