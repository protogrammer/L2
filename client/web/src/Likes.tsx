import {AiFillHeart} from 'react-icons/ai'
import {MsgData, toggleMessageLike} from './Message'

type Props = {
    item: MsgData
    data: MsgData[]
    setData: React.Dispatch<React.SetStateAction<MsgData[]>>
}

export const Likes: React.FC<Props> = (props: Props) => {
    const style = props.item.liked ? {color: 'red'} : {color: 'lightgray'};
    return (<>
        <AiFillHeart style={style} size={20} onClick={() => {
            toggleMessageLike(props.item, props.data, props.setData)
        }}/>
        {props.item.likes}
    </>);
}