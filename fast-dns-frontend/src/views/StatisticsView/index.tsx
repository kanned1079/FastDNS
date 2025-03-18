import React, {useMemo} from "react";
import {Card, Button, Typography} from "antd";
import styles from "../Dashboard/index.module.less";

type StatisticProps = {
    cache_rate: number
    cache_hint: number
    cache_miss: number
    updateData: () => void
}

const {Title} = Typography

const StatisticsView: React.FC<StatisticProps> = ({cache_rate, cache_hint, cache_miss, updateData}) => {
    const reqAll = useMemo(() => {
        if (cache_miss === undefined || cache_hint === undefined)
            return 0
        return cache_miss + cache_hint
    }, [cache_miss, cache_hint])

    const topCardContent = [
        {
            title: 'DNS缓存命中率',
            data: cache_rate,
            type: 'rate',
            description: '如果缓存中已有该域名的解析记录，可以直接返回该结果，避免再次向上游 DNS 服务器发起请求。',
        },
        {
            title: 'DNS缓存命中次数',
            data: cache_hint,
            type: 'number',
            description: '具体来说，缓存提示（hint）通常表示该记录在缓存中，但它可能存在一些过期、失效或者需要再次验证的风险。',

        },
        {
            title: '没有找到缓存',
            data: cache_miss,
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

    return (
        <div className={styles.root}>
            <div
                style={{display: 'flex', justifyContent: 'flex-start', alignItems: 'center', marginBottom: '12px'} as React.CSSProperties}
            >
                <Title style={{marginBottom: 0}} level={3}>统计数据</Title>
                <Button style={{marginLeft: '10px'}} color={"primary"} variant={"solid"} size={'small'} onClick={updateData}>刷新统计数据</Button>
            </div>

            <div style={{
                display: 'flex',
                flexDirection: 'row',
                justifyContent: 'center',
                width: '100%',
                gap: '20px'
            } as React.CSSProperties}>
                {topCardContent.map((item) => (
                    <Card variant={"borderless"} className={styles.card} key={item.title}>
                        <h4 style={{opacity: 0.6} as React.CSSProperties}>{item.title}</h4>
                        <p style={{
                            fontSize: '1.3rem',
                            fontWeight: 'bold'
                        } as React.CSSProperties}>{item.type === "rate" && item.data
                            ? (item.data * 100).toFixed(2) + " %"
                            : item.data + ' 次'
                        }</p>
                        <p style={{marginTop: '10px', opacity: '0.8'} as React.CSSProperties}>{item.description}</p>
                    </Card>
                ))}
            </div></div>
    )

}

export default StatisticsView