import React from "react";
import styles from "./index.module.less"

// 声明 Props 类型
type Props = {
    title: string;
    description?: string;
};

// 正确声明组件类型
const PageHead: React.FC<Props> = ({ title, description }) => {
    return (
        <div className={styles.root}>
            <p className={styles.title}>{title}</p>
            <p className={styles.desc}>{description}</p>
        </div>
    );
};

export default PageHead;