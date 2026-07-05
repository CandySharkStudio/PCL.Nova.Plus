<script lang="ts">
    import MySelectCard from "../../../component/card/MySelectCard.svelte";
    import MyNormalButton from "../../../component/button/MyNormalButton.svelte";
    import {
        inputbox,
        messagebox,
        showHint,
        MSG_ERROR,
        MSG_INFO,
        MSG_WARNING,
    } from "../../../store/messagebox";
    import {
        dont_click,
        download_list,
        unlock_theme,
    } from "../../../store/changeBody";
    import {
        GetSelfUniqueAddress,
        OpenCustomURL,
        sleep,
    } from "../../../store/functions";
    import {
        GetPCL2Identify,
        CryptoUnlock,
    } from "../../../../wailsjs/go/main/App";
    import { GetTodayLucky } from "../../../../wailsjs/go/main/App.js";
    import MyNormalSpan from "../../../component/input/MyNormalSpan.svelte";
    import MyTextButton from "../../../component/button/MyTextButton.svelte";
    import MyTextInput from "../../../component/input/MyTextInput.svelte";
    import {
        WriteConfig,
        ReadConfig,
        OpenDirectoryDialog,
        OpenExplorer,
    } from "../../../../wailsjs/go/launcher/ReaderWriter";
    import { onMount } from "svelte";

    export let slide = null;
    export let after_leave = null;

    let savePath = "";
    let headers = "";
    let cookies = "";
    let url = "";
    let fileName = "";
    function getFileNameFromURL(url: string) {
        return url.substring(url.lastIndexOf("/") + 1);
    }
    function changeUrl(e: CustomEvent) {
        url = e.detail.value;
        fileName = getFileNameFromURL(url);
    }
    async function copyNovaIdentify() {
        let add = await GetSelfUniqueAddress();
        if (add != "") {
            navigator.clipboard.writeText(add);
            showHint("唯一识别码已复制到剪贴板");
        }
    }
    async function changeSavePath(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig("current", "Download", "SavePath", value.toString());
        savePath = value;
    }
    async function changeHeaders(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig("current", "Download", "Headers", value.toString());
        headers = value;
    }
    async function changeCookies(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig("current", "Download", "Cookies", value.toString());
        cookies = value;
    }
    async function selectFolder() {
        let path = await OpenDirectoryDialog("请选择保存路径！");
        if (path == "") {
            return;
        }
        await WriteConfig("current", "Download", "SavePath", path.toString());
        savePath = path;
    }
    onMount(async () => {
        headers = await ReadConfig("current", "Download", "Headers");
        savePath = await ReadConfig("current", "Download", "SavePath");
        cookies = await ReadConfig("current", "Download", "Cookies");
    });
    async function dontClick() {
        await messagebox(
            "警告",
            "PCL Nova Plus 作者不会受理由于点击千万别点造成的任何 Bug。这是最后的警告，是否继续操作？",
            MSG_ERROR,
            ["ok", "ok", "ok"],
        );
        let rand = Math.floor(Math.random() * 3 + 1);
        // let rand = 2
        dont_click.set(rand);
    }
    async function todayLucky() {
        let luck = await GetTodayLucky(await GetPCL2Identify());
        if (luck == -1) {
            await messagebox(
                "警告",
                "PCL Nova Plus 无法正常获取到 当前硬件配置UUID 你可能在虚拟机中运行 Nova Plus!",
                MSG_WARNING,
            );
        } else if (luck == -2) {
            await messagebox(
                "警告",
                "PCL Nova Plus 无法写入 PCL Identify，请尝试用管理员身份启动一次，随后再使用 Nova 查看今日人品~",
                MSG_WARNING,
            );
        } else if (luck == -3) {
            let i = await inputbox(
                "被你发现了！",
                "PCL Nova Plus 虽然无法为除了 Windows 以外的操作系统提供今日人品，但是龙猫许可证规定了 Windows 不能有主题，而别的系统可以有呀~因此请在下面输入解锁码，可以解锁主题蒙版噢~<br><br>解锁码获取条件：在 PCL.Nova.Plus 仓库里提交任意一个<b>有建设性的</b> Pull Request 或者任意一个<b>有建设性的</b> Issue，随后将你的唯一识别码放在下面，之后作者会在底下回复你的解锁码嗷~",
            );
            if (i == "") return;
            let sua = await GetSelfUniqueAddress();
            if (sua == "") return;
            if (await CryptoUnlock(i)) {
                await messagebox(
                    "解锁成功！",
                    "恭喜解锁成功主题蒙版！现在你可以正常使用 Nova Plus 的主题啦！",
                );
                await WriteConfig("other", "Unlock", sua, i);
                unlock_theme.set(true);
            } else {
                await messagebox(
                    "解锁失败！",
                    "解锁码格式错误，请重试！",
                    MSG_ERROR,
                );
            }
        } else {
            let str = "";
            if (luck == 100) {
                str = "！100！100！！！！！";
            } else if (luck == 99) {
                str = "！额，但不是 100……";
            } else if (luck >= 90) {
                str = "！好评如潮！";
            } else if (luck >= 60) {
                str = "！是不错的一天呢！";
            } else if (luck > 50) {
                str = "！还行啦还行啦。";
            } else if (luck == 50) {
                str = "！五五开……";
            } else if (luck >= 40) {
                str = "！还……还行吧……？";
            } else if (luck >= 11) {
                str = "！呜哇……";
            } else if (luck >= 1) {
                str = "……（没错，是百分制）";
            } else {
                str = "……";
                if (
                    (await messagebox(
                        "今日人品 - 附加使用条款",
                        "在查看结果前，请先同意以下附加使用条款：<br><br>1. 我知晓并了解 PCL.Nova.Plus 的今日人品功能完全没有出 Bug。<br>2. PCL.Nova.Plus 不对使用本软件所间接造成的一切财产损失（如砸电脑等）等负责。",
                        MSG_ERROR,
                        ["再见", "同意并查看结果"],
                    )) == 0
                ) {
                    return;
                }
            }
            const date = new Date();
            await messagebox(
                "今日人品 - " +
                    date.getFullYear() +
                    "/" +
                    (date.getMonth() + 1) +
                    "/" +
                    date.getDate(),
                "你今天的人品值是：" + luck + str,
                MSG_INFO,
                ["我知道了"],
            );
        }
    }
    let myCaveDOM: HTMLDivElement | null = null;
    let curChange = false;
    let cave =
        "反复点击这里可以查看 PCL.Nova 作者与各位沙雕网友乱七八糟的留言！";
    import monishuju from "../../../store/cave.json";
    async function showCave() {
        let f = false;
        for (let i = 10; i >= 1; i--) {
            myCaveDOM!.style.opacity = f ? "1" : "0";
            await sleep(i * 10);
            f = !f;
        }
        let tCave = monishuju[Math.floor(Math.random() * monishuju.length)];
        cave = "";
        for (let i = 0; i < tCave.length; i++) {
            cave += tCave[i];
            await sleep(20);
        }
        curChange = false;
    }
    async function openDirectory() {
        if (savePath == "") return;
        await OpenExplorer(savePath);
    }
    async function startDownload() {
        download_list.update((list) => {
            list.push({
                savePath: savePath + "/" + fileName,
                url: url,
            });
            return list;
        });
    }
    function caveClick() {
        if (curChange) return;
        curChange = true;
        showCave();
    }
</script>

<div
    class="component-treasure_box"
    in:slide
    out:slide
    on:outroend={after_leave}
>
    <MySelectCard title="百宝箱">
        <div class="version-all">
            <div style="width: 100%">
                {#if $dont_click === 0}
                    <MyNormalButton
                        style_in="border: 1px solid red; color: red; padding: 0 20px; height: 35px"
                        on:click={dontClick}>千万别点</MyNormalButton
                    >
                {/if}
                <MyNormalButton
                    style_in="padding: 0 20px; height: 35px; margin-left: 10px"
                    on:click={todayLucky}>今日人品</MyNormalButton
                >
                <MyNormalButton
                    style_in="padding: 0 20px; height: 35px; margin-left: 10px"
                    on:click={copyNovaIdentify}>复制 Nova 识别码</MyNormalButton
                >
            </div>
        </div>
    </MySelectCard>
    <MySelectCard title="">
        <div class="title-bar">
            <MyNormalSpan style_in="font-size: 1.17em; font-weight: bold;">
                回声洞
            </MyNormalSpan>
            <MyTextButton
                on:click={() =>
                    OpenCustomURL("https://wj.qq.com/s2/24355863/9e54/")}
                style_in="margin: 0 10px;"
                hover_style_in="margin: 0 10px; color: #20ACFF"
                active_style_in="margin: 0 10px; color: #008CFF"
                >投稿</MyTextButton
            >
        </div>
        <div class="version-all">
            <div
                on:click={caveClick}
                on:keydown|preventDefault
                class="cave-div"
            >
                <div bind:this={myCaveDOM}>
                    {@html cave}
                </div>
            </div>
        </div>
    </MySelectCard>
    <MySelectCard title="下载自定义文件（暂未完成）">
        <div class="version-all">
            <MyNormalSpan
                title={"但是各位可以自定义 header 甚至是 Cookie！头，可以在 headers 里自定义 Authorization 或者 Referer！这样应该可以实现 PCL 高速下载了💦\n不过建议各位还是保持默认的 header、Cookie 吧~\n在 PCL.Nova.Plus 里，自定义下载的 User-Agent 默认保持的是 Chrome 浏览器的，目前暂时无法自定义~\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.76\nContent-Type 默认使用的是 application/octet-stream。目前暂时无法自定义~\n上述两个会覆盖你们所写的内容。。"}
            >
                使用 PCL.Nova.Plus
                的高速多线程下载引擎下载任意文件。请注意，部分网站（例如百度网盘）可能会报错（403）已禁止，无法正常下载。
            </MyNormalSpan>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >下载地址</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    placeholder="请输入下载地址"
                    on:blur={changeUrl}
                />
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >保存到</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    placeholder="请输入保存路径"
                    on:blur={changeSavePath}
                    value={savePath}
                />
                <MyTextButton
                    on:click={selectFolder}
                    style_in="margin: 0 10px;"
                    hover_style_in="margin: 0 10px; color: #20ACFF"
                    active_style_in="margin: 0 10px; color: #008CFF"
                    >选择</MyTextButton
                >
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >文件名</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    placeholder="请输入文件名"
                    value={fileName}
                />
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >headers</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    placeholder="请输入 Headers（请保留原始 JSON 形式）（如果不确定请不要填）"
                    title="请保留原始 JSON 形式。"
                    on:blur={changeHeaders}
                    value={headers}
                />
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >Cookie</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    placeholder="请输入 Cookie（请保留【key1=value1; key2=value2;】这种形式的）（如果不确定请不要填）"
                    title="请保留【key1=value1; key2=value2;】这种形式的。"
                    on:blur={changeCookies}
                    value={cookies}
                />
            </div>
            <div class="proc">
                <div style="display: flex;">
                    <MyNormalButton
                        style_in="width: 130px; height: 35px; margin-right: 20px"
                        on:click={startDownload}>开始下载</MyNormalButton
                    >
                    <MyNormalButton
                        style_in="width: 130px; height: 35px;"
                        on:click={openDirectory}>打开文件夹</MyNormalButton
                    >
                </div>
            </div>
        </div>
    </MySelectCard>
</div>

<style>
    .component-treasure_box {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        overflow-y: auto;
    }
    .title-bar {
        width: calc(100% - 40px);
        display: flex;
        padding: 10px 20px;
        align-items: center;
        justify-content: space-between;
    }
    .cave-div {
        width: 100%;
        word-wrap: break-word;
        white-space: pre-wrap;
    }
    .settings {
        width: 100%;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
</style>
