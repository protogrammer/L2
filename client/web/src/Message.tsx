import React from 'react'
import {MessageButtons} from './MessageButtons'
import {Likes} from './Likes'

export type MsgData = {
    id: number
    created: Date
    author: string | null
    text: string
    likes: number
    liked: boolean
}

const srvGetMessages = (): Promise<MsgData[] | null> => {
    return fetch("/api/get").then(data => data.json())
}

const srvNewMessage = (text: string): Promise<MsgData | null> => {
    return fetch(`/api/new`, {
        method: 'POST',
        body: text,
    }).then(data => data.json())
}

const srvUpdateMessage = (id: number, text: string): Promise<string | null> => {
    return fetch(`/api/edit?id=${id}`, {
        method: 'POST',
        body: text,
    }).then(() => text)
}

const srvDeleteMessage = (id: number): Promise<boolean> => {
    return fetch(`/api/delete?id=${id}`, {
        method: "POST",
    }).then(() => true)
}

const srvToggleMessageLike = (id: number, value: boolean): Promise<boolean> => {
    return fetch(`/api/like?id=${id}&value=${value}`, {
        method: 'POST'
    }).then(data => Number(data) != null)
}

export const getMessages = (setData: React.Dispatch<React.SetStateAction<MsgData[]>>) => {
    console.log('getMessages: ');
    srvGetMessages().then((list) => {
        if (!list) {
            alert('can not get posts');
            return;
        }
        setData(list);
    });
}

export const editMessage = (msg: MsgData, data: MsgData[], setData: React.Dispatch<React.SetStateAction<MsgData[]>>) => {
    console.log('editMessage: ', msg.id);
    const result = prompt("Please change your message", msg.text);
    if (result && result !== msg.text) {
        {/* update message on server */
        }
        srvUpdateMessage(msg.id, result).then((text) => {
            if (text === null) {
                alert('server error occurred');
                return
            }
            const editedMsg: MsgData = {...msg, text: text}
            setData(data.map((otherMsg) => otherMsg.id !== msg.id ? otherMsg : editedMsg))
        })
        // const msgtextdiv = document.getElementById( 'msgtext'+msg.id );
        // if( msgtextdiv ) msgtextdiv.innerText = msg.text;
    }
}

export const deleteMessage = (msg_id: number, data: MsgData[], setData: React.Dispatch<React.SetStateAction<MsgData[]>>) => {
    console.log('deleteMessage: ', msg_id);
    srvDeleteMessage(msg_id).then((ok) => {
        if (ok) setData(data.filter((otherMsg) => otherMsg.id !== msg_id))
    });
}

export const newMessage = (data: MsgData[], setData: React.Dispatch<React.SetStateAction<MsgData[]>>) => {
    console.log('newMessage');
    const result = prompt("Please write your message", '');
    if (result && result.length > 0) {
        srvNewMessage(result).then((msg) => {
            console.log(msg);
            if (!msg) {
                alert('Error saving message');
                return;
            }
            setData([msg, ...data]);
        });
    }

}

export const toggleMessageLike = (msg: MsgData, data: MsgData[], setData: React.Dispatch<React.SetStateAction<MsgData[]>>) => {
    console.log('toggleMessageLike: ', msg.id);
    srvToggleMessageLike(msg.id, !msg.liked).then((ok) => {
        if (ok) {
            const editedMsg: MsgData = {...msg, liked: !msg.liked}
            if (editedMsg.liked) {
                editedMsg.likes++
            } else {
                editedMsg.likes--
            }
            setData(data.map((otherMsg) => otherMsg.id !== msg.id ? otherMsg : editedMsg))
        }
    });
}

type Props = {
    item: MsgData
    data: MsgData[]
    setData: React.Dispatch<React.SetStateAction<MsgData[]>>
}

export const Message: React.FC<Props> = (props: Props) => {
    // console.log( 'Message.render', props.item.id );
    return <div className={[
        'message',
        `${props.item.author === null ? 'mine' : ''}`
    ].join(' ')}>
        {props.item.text}
        {props.item.author === null
            ? <div className="status">
                <Likes item={props.item} data={props.data} setData={props.setData}/>&emsp;
                <MessageButtons msg={props.item} data={props.data} setData={props.setData}/>
            </div>
            : <div className="status" style={{color: 'blue'}}>
                <Likes item={props.item} data={props.data} setData={props.setData}/>&emsp;
                {props.item.author}&nbsp;{props.item.created.toLocaleString()}
            </div>}

    </div>
}