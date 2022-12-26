import React, {useEffect, useState} from 'react'
import './index.css';
import {Message, MsgData, newMessage, getMessages} from './Message'
import {BiMessageAdd} from "react-icons/bi";

export const FeedScreen: React.FC = () => {
    const [data, setData] = useState<MsgData[]>([]);
    useEffect(() => {
        getMessages(setData)
    }, [])
    const messages = data.map((item) => {
        return <Message item={item} data={data} setData={setData}/>
    });
    return (
        <ul>
            {messages}
            <div className="message" style={{background: 0xffffff, alignItems: 'flex-end'}}>
                <button style={{background: 0xffffff, color: 'blue'}} onClick={() => {
                    newMessage(data, setData)
                }}><BiMessageAdd size={50}/></button>
            </div>
        </ul>
    );
}