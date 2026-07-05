<script lang="ts">
    import MySelectCard from "../../../component/card/MySelectCard.svelte";
    import MyNormalSpan from "../../../component/input/MyNormalSpan.svelte";
    import MyProgressBar from "../../../component/input/MyProgressBar.svelte";
    import MyRadioButton from "../../../component/button/MyRadioButton.svelte";
    import MyTextInput from "../../../component/input/MyTextInput.svelte";
    import {
        ReadConfig,
        WriteConfig,
    } from "../../../../wailsjs/go/launcher/ReaderWriter";
    import { onMount } from "svelte";

    export let slide = null;
    export let after_leave = null;
    let maximumThread = 64;
    let downloadSource = 1;
    let proxyUrl = "";
    let proxyPort = "";
    let proxyUsername = "";
    let proxyPassword = "";
    let proxyType = 1;
    async function changeDownloadSource(source: number) {
        await WriteConfig(
            "current",
            "Version",
            "SelectDownloadSource",
            source.toString(),
        );
        downloadSource = source;
    }
    async function changeMaximumThread(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig(
            "current",
            "Version",
            "ThreadBiggest",
            value.toString(),
        );
        maximumThread = value;
    }
    async function changeProxyUrl(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig("current", "Download", "ProxyUrl", value.toString());
        proxyUrl = value;
    }
    async function changeProxyPort(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig("current", "Download", "ProxyPort", value.toString());
        proxyPort = value;
    }
    async function changeProxyUsername(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig(
            "current",
            "Download",
            "ProxyUsername",
            value.toString(),
        );
        proxyUsername = value;
    }
    async function changeProxyPassword(e: CustomEvent) {
        let value = e.detail.value;
        await WriteConfig(
            "current",
            "Download",
            "ProxyPassword",
            value.toString(),
        );
        proxyPassword = value;
    }
    async function changeProxyType(source: number) {
        await WriteConfig(
            "current",
            "Download",
            "ProxyType",
            source.toString(),
        );
        proxyType = source;
    }
    onMount(async () => {
        let dSrc = parseInt(
            await ReadConfig("current", "Version", "SelectDownloadSource"),
        );
        if (!Number.isNaN(dSrc) && dSrc > 0 && dSrc < 4) {
            downloadSource = dSrc;
        }
        let mThr = parseInt(
            await ReadConfig("current", "Version", "ThreadBiggest"),
        );
        if (!Number.isNaN(mThr) && mThr > 0 && mThr <= 256) {
            maximumThread = mThr;
        }
        let pType = parseInt(
            await ReadConfig("current", "Download", "ProxyType"),
        );
        if (!Number.isNaN(pType) && pType > 0 && mThr <= 3) {
            proxyType = pType;
        }
        proxyUrl = await ReadConfig("current", "Download", "ProxyUrl");
        proxyPort = await ReadConfig("current", "Download", "ProxyPort");
        proxyUsername = await ReadConfig(
            "current",
            "Download",
            "ProxyUsername",
        );
        proxyPassword = await ReadConfig(
            "current",
            "Download",
            "ProxyPassword",
        );
    });
</script>

<div class="component-other" in:slide out:slide on:outroend={after_leave}>
    <MySelectCard title="下载">
        <div class="version-all">
            <div class="settings">
                <MyNormalSpan style_in="width: 120px">最大线程数：</MyNormalSpan
                >
                <div class="setting-con">
                    <MyProgressBar
                        style_in="width: calc(100% - 50px);"
                        max={256}
                        min={1}
                        on:dragging={changeMaximumThread}
                        value={maximumThread}
                    />
                    <MyNormalSpan>{maximumThread}</MyNormalSpan>
                </div>
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 120px">下载源选择：</MyNormalSpan
                >
                <div class="setting-con">
                    <MyRadioButton
                        on:click={() => changeDownloadSource(1)}
                        isChecked={downloadSource === 1}>官方源</MyRadioButton
                    >
                    <MyRadioButton
                        on:click={() => changeDownloadSource(2)}
                        isChecked={downloadSource === 2}>BMCLAPI</MyRadioButton
                    >
                    <MyRadioButton
                        on:click={() => changeDownloadSource(3)}
                        isChecked={downloadSource === 3}>二者平衡</MyRadioButton
                    >
                </div>
            </div>
        </div>
    </MySelectCard>
    <MySelectCard title="高级设置" canExpand isExpand>
        <div class="version-all">
            <div class="red-hint hint" style="width: calc(100% - 30px)">
                请谨慎修改以下设置，除非你真的确定你在做什么
            </div>
            <h3><MyNormalSpan>代理设置</MyNormalSpan></h3>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >代理类型：</MyNormalSpan
                >
                <div class="setting-con" style="margin-left: 20px;">
                    <MyRadioButton
                        on:click={() => changeProxyType(1)}
                        isChecked={proxyType === 1}>HTTP</MyRadioButton
                    >
                    <MyRadioButton
                        on:click={() => changeProxyType(2)}
                        isChecked={proxyType === 2}>HTTPS</MyRadioButton
                    >
                    <MyRadioButton
                        on:click={() => changeProxyType(3)}
                        isChecked={proxyType === 3}>SOCKS5</MyRadioButton
                    >
                </div>
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >主机网址：</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    title="只需要输入地址即可，无需输入 https、http、sock5等"
                    placeholder="请输入主机地址"
                    on:blur={changeProxyUrl}
                    value={proxyUrl}
                />
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >主机端口：</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    title="代理一般都是有端口的，填入端口即可"
                    placeholder="请输入主机端口（如果没有可以填80）"
                    on:blur={changeProxyPort}
                    value={proxyPort}
                />
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >账号：</MyNormalSpan
                >
                <MyTextInput
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    title="如果需要身份验证，才填，否则不需要填"
                    placeholder="如果你的代理需要身份验证，则这里填账号，否则请留空"
                    on:blur={changeProxyUsername}
                    value={proxyUsername}
                />
            </div>
            <div class="settings">
                <MyNormalSpan style_in="width: 80px; text-align: right;"
                    >密码：</MyNormalSpan
                >
                <MyTextInput
                    password
                    style_in="flex: 1; margin-left: 20px; height: 25px"
                    title="如果需要身份验证，才填，否则不需要填"
                    placeholder="如果你的代理需要身份验证，则这里填密码，否则请留空"
                    on:blur={changeProxyPassword}
                    value={proxyPassword}
                />
            </div>
        </div>
    </MySelectCard>
</div>

<style>
    .component-other {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        overflow-y: auto;
    }
    .settings {
        width: 100%;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    .setting-con {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
</style>
