// import {Button} from "antd"
import {Layout, ConfigProvider, theme} from 'antd';
import type {ThemeConfig} from 'antd'
// import styles from "./App.module.less"
import Dashboard from "./views/Dashboard";

const {Header, Footer, Content} = Layout;

const headerStyle: React.CSSProperties = {
    width: '100%',
    height: '64px',
    paddingLeft: '20px',
}

const contentStyle: React.CSSProperties = {
    width: '100%',
    height: 'calc(100vh - 128px)',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    paddingTop: '30px'
}

const customTheme: ThemeConfig = {
    token: {
        colorPrimary: '#6390b9',

    },
    components: {
        Card: {
            colorBgContainer: '#fff',
        },
        Layout: {
            bodyBg: '#eaeaea',
            headerBg: '#6390b9'
        }
    },
    algorithm: theme.defaultAlgorithm
}


function App() {
    return (
        <ConfigProvider
            theme={customTheme}>
            <Layout>
                <Header style={headerStyle}>
                    <p style={{fontSize: '1.2rem', color: '#fff'} as React.CSSProperties}>
                        FastDNS Dashboard
                    </p>
                </Header>
                <Content style={contentStyle}>
                    <Dashboard/>
                </Content>
                <Footer style={headerStyle}>Footer</Footer>
            </Layout>
        </ConfigProvider>

    )
}

export default App