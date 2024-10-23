import React from 'react';
import { Segment, Dimmer, Loader } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css';

const Loading = ({ prompt: name = 'page' }) => {
  return (
    <Segment       
      style={{         
        minHeight: '100vh',        // 高度设置为视口高度
        display: 'flex',           // 使用 flex 布局
        justifyContent: 'center',  // 水平居中
        alignItems: 'center',      // 垂直居中
      }}
    >
      <Dimmer active inverted>
        <Loader indeterminate>加载{name}中...</Loader>
      </Dimmer>
    </Segment>
  );
};

export default Loading;
