import instance from "../../axios/index"
import React, {useEffect, useState} from "react";

import PageHead from "../PageHead"
import {type Response} from "../../types.tsx";

import styles from "./index.module.less"


import StatisticsView from "../StatisticsView";
import {Tag} from 'antd'


function Dashboard() {
    const [data, setData] = useState<Response>()

    const fetchDataClick = async () => {
        console.log('获取数据')
        let {data} = await instance.get('/api/statistics')
        if (data.code === 200) {
            console.log(data.config)
            setData(data)
        }
    }
    // fetchDataClick()
    let intervalId: number|undefined = undefined

    const dnsConfig: {
        title: string
        data: any
        type: 'list' | 'bool' | 'string'
    }[] = [
        {
            title: '主要DNS列表',
            data: data?.config.dns.list,
            type: 'list',
        },
        {
            title: '备用DNS列表',
            data: data?.config.dns.backup_secondary_list,
            type: 'list'
        },
        {
            title: '请求类型',
            data: data?.config.dns.mode,
            type: 'string'
        },
        {
            title: '是否启用带有tls的DNS',
            data: data?.config.dns.tls_enabled,
            type: 'bool'
        }
    ]

    useEffect(() => {
        console.log('dashboard mount')
        fetchDataClick()
        // startLookup()
        return () => {
            intervalId?clearInterval(intervalId):null
        }
    }, []);

    return (
        <div className={styles.root}>
            {/*<PageHead*/}
            {/*    title={'仪表板'}*/}
            {/*/>*/}

            <StatisticsView
                cache_rate={data?.cache.cache_rate || 0}
                cache_hint={data?.cache.cache_hint || 0}
                cache_miss={data?.cache.cache_miss || 0}
                updateData={fetchDataClick}
            ></StatisticsView>

            <PageHead
                title={'服务器配置'}
                description={'在这里将显示具体的配置信息，如使用的DNS服务器主从列表、是否启用EDNS和SCDNS等高级应用'}
            ></PageHead>

            {dnsConfig.map((item) => (
                <div
                    key={item.title} style={{display: 'flex', marginBottom: '10px'} as React.CSSProperties}>
                    <p style={{fontSize: '1rem', fontWeight: 'bold', marginRight: '10px'} as React.CSSProperties}>{item.title}</p>
                    {item.type === 'list' && Array.isArray(item.data) && item.data.length
                        ? item.data.map((str: string, index) => (
                            <Tag color={'cyan'} key={index}>{str}</Tag>
                        ))
                        : item.type === 'bool'
                            ? (item.data ?(<Tag style={{padding: '0 14px'}} color={'green'}>{'是'}</Tag>) : (<Tag color={'error'}>否</Tag>))
                            : (<div>{item.data}</div>)
                    }
                </div>
            ))}

        </div>
    )
}

export default Dashboard