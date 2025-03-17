// import {Button} from "antd"
import { Layout } from 'antd';
// import styles from "./App.module.less"
import Dashboard from "./views/Dashboard";

const { Header, Footer, Content } = Layout;

const headerStyle: React.CSSProperties = {
    width: '100%',
    height: '64px',
    backgroundColor: '#5e95cd',
    paddingLeft: '20px',
}

const contentStyle: React.CSSProperties = {
    width: '100%',
    height: 'calc(100vh - 128px)',
    backgroundColor: '#e3e5e7'
}



function App() {

    return (
        <Layout>
            <Header style={headerStyle}>
                <p style={{fontSize: '1.2rem', color: '#fff'} as React.CSSProperties}>
                    FastDNS Dashboard
                </p>
            </Header>
            <Content style={contentStyle}>
                <Dashboard />
            </Content>
            <Footer style={headerStyle}>Footer</Footer>
        </Layout>
    )
}

export default App