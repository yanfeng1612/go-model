module.exports = {
    // baseUrl: './',
    assetsDir: 'static',
    productionSourceMap: false,
    devServer: {
        proxy: {
            '/api':{
                // target:'http://jsonplaceholder.typicode.com',
                target:'http://localhost:1601/',
                changeOrigin:true,
                pathRewrite:{
                    '/api':''
                }
            }
        }
    }
}