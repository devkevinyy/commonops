const { override, fixBabelImports } = require('customize-cra');
// const { addLessLoader } = require('customize-cra');

module.exports = override(
    fixBabelImports('import', {
        libraryName: 'antd',
        libraryDirectory: 'es',
        style: 'css',
    }),
    // addLessLoader({
    //     javascriptEnabled: true,
    //     modifyVars: { '@primary-color': '#ff6a00' },
    // }),
);