export default {
  module: {
    rules: [
      {
        test: /\.css$/,
        loader: 'postcss-loader',
        options: {
          ident: 'postcss',
          plugins: () => [
            require('postcss-import'),
            require('tailwindcss')('./tailwind.config.ts'),
            require('autoprefixer'),
          ],
        },
      },
    ],
  },
};
