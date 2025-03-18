import instance from "../../axios/index"
import {useEffect, useState, useMemo} from "react";
import {Card} from "antd"
import PageHead from "../PageHead"
import {type Response} from "../../types.tsx";

import styles from "./index.module.less"
import * as React from "react";

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
    const startLookup = () => {
        intervalId = setInterval(() => fetchDataClick(), 1000)
    }

    const reqAll = useMemo(() => {
        if (data?.cache.cache_miss === undefined || data?.cache.cache_hint === undefined)
            return 0
        return data?.cache.cache_miss + data?.cache.cache_hint
    }, [data?.cache.cache_miss, data?.cache.cache_hint])

    const topCardContent = [
        {
            title: 'DNS缓存命中率',
            data: data?.cache.cache_rate,
            type: 'rate',
            description: '如果缓存中已有该域名的解析记录，可以直接返回该结果，避免再次向上游 DNS 服务器发起请求。',
        },
        {
            title: 'DNS缓存命中次数',
            data: data?.cache.cache_hint,
            type: 'number',
            description: '具体来说，缓存提示（hint）通常表示该记录在缓存中，但它可能存在一些过期、失效或者需要再次验证的风险。',

        },
        {
            title: '没有找到缓存',
            data: data?.cache.cache_miss,
            type: 'number',
            description: '发生在系统无法从缓存中获取到所请求的域名解析记录时。此时，DNS 系统会向上游的 DNS 服务器发起请求，从而获取新的解析结果。',
        },
        {
            title: '请求总计',
            data: reqAll,
            type: 'number',
            description: '包含成功和失败相应的计数。',
        },
    ]

    useEffect(() => {
        console.log('dashboard mount')
        fetchDataClick()
        startLookup()
        return () => {
            intervalId?clearInterval(intervalId):null
        }
    }, []);

    return (
        <div className={styles.root}>
            <PageHead
                title={'DNS记录概览'}
                description={'在这里展示的是DNS服务器的概览信息'}
            ></PageHead>
            <div style={{
                display: 'flex',
                flexDirection: 'row',
                justifyContent: 'center',
                width: '100%',
                padding: '20px',
                gap: '20px'
            } as React.CSSProperties}>
                {topCardContent.map((item) => (
                    <Card variant={"borderless"} className={styles.card} key={item.title}>
                        <h4 style={{opacity: 0.8} as React.CSSProperties}>{item.title}</h4>
                        <p style={{
                            fontSize: '1.3rem',
                            fontWeight: 'bold'
                        } as React.CSSProperties}>{item.type === "rate" && item.data
                            ? (item.data * 100).toFixed(2) + " %"
                            : item.data + ' 次'
                        }</p>
                        {/*<p style={{marginTop: '10px', opacity: '0.8'} as React.CSSProperties}>{item.description}</p>*/}
                    </Card>
                ))}
            </div>
            <PageHead
                title={'服务器配置'}
                description={'在这里将显示具体的配置信息，如使用的DNS服务器主从列表、是否启用EDNS和SCDNS等高级应用'}
            ></PageHead>



        </div>
    )
}

export default Dashboard