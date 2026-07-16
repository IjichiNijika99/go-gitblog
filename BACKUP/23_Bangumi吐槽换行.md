# [Bangumi吐槽换行](https://github.com/IjichiNijika99/go-gitblog/issues/23)

### [吐槽简评换行显示](https://bgm.tv/dev/app/6622) / v1.0.2

给条目页、收藏列表、时间胶囊及收藏盒中的“吐槽/简评”添加回车换行显示，为强迫症用户还原真实排版。

```
// ==UserScript==
// @name         Bangumi 吐槽换行显示
// @namespace    https://bgm.tv/
// @version      1.2
// @description  我要看 Bangumi 吐槽换行。隔离条目页与收藏页排版样式，防止排版错乱
// @author       Gumiku
// @match        *://bgm.tv/*
// @match        *://bangumi.tv/*
// @match        *://chii.in/*
// @grant        GM_addStyle
// @run-at       document-end
// ==/UserScript==

(function() {
    'use strict';

    // 注入 CSS：区分条目页的 .comment 和收藏列表页的 .text
    const css = `
        /* 绝大多数场景（条目页吐槽箱、各类主页、吐槽页面）：文字都在 .comment 标签里 */
        .comment,

        /* 仅限个人收藏/列表/目录页（/collect、/index等）：限定在 #browserItemList 内部的 .text，避免污染条目页布局 */
        #browserItemList .comment-box .text,

        /* 时光机与时间胶囊的文字 */
        #timeline .text,
        #timeline .status,

        /* 右侧边栏收藏盒（脚本自己生成的标签） */
        .my-side-comment {
            white-space: pre-wrap !important;
            word-break: break-word !important;
        }
    `;

    if (typeof GM_addStyle !== 'undefined') {
        GM_addStyle(css);
    } else {
        const style = document.createElement('style');
        style.textContent = css;
        document.head.appendChild(style);
    }

    // 修复右侧边栏「收藏盒」中的裸文本节点吐槽
    function fixSidePanelComment() {
        const sidePanel = document.querySelector('#panelInterestWrapper .SidePanel');
        if (!sidePanel || sidePanel.querySelector('.my-side-comment')) return;

        sidePanel.childNodes.forEach(node => {
            if (node.nodeType === Node.TEXT_NODE && node.textContent.trim().length > 0) {
                const wrapper = document.createElement('div');
                wrapper.className = 'my-side-comment';
                wrapper.style.cssText = 'padding: 4px 0; line-height: 1.5;';
                wrapper.textContent = node.textContent.trim();
                node.replaceWith(wrapper);
            }
        });
    }

    // 清除后端模板自带的首行多余空格/换行（不破坏表情和链接标签）
    function cleanTemplateWhitespace() {
        const selectors = [
            '.comment',
            '#browserItemList .comment-box .text',
            '.timeline_img .text'
        ].join(',');

        document.querySelectorAll(selectors).forEach(el => {
            // 仅安全修改第一个 DOM 文本节点，消除 ^\s+ (开头所有空白和换行)，防止覆盖 <a> 链接或 <img> 表情
            if (el.firstChild && el.firstChild.nodeType === Node.TEXT_NODE) {
                el.firstChild.nodeValue = el.firstChild.nodeValue.replace(/^\s+/, '');
            }
        });
    }

    // 初始化
    function init() {
        fixSidePanelComment();
        cleanTemplateWhitespace();
    }

    init();

    // 监听 DOM 变化：兼容 AJAX 局部无刷新修改
    const observer = new MutationObserver(() => {
        init();
    });

    observer.observe(document.body, {
        childList: true,
        subtree: true
    });
})();
```

### 效果图

<img width="749" height="254" alt="Image" src="https://github.com/user-attachments/assets/63d50dce-500f-47e8-8725-5b9c526fb35a" />

<img width="735" height="305" alt="Image" src="https://github.com/user-attachments/assets/8b978533-168c-49e9-af9a-ca61f7728131" />