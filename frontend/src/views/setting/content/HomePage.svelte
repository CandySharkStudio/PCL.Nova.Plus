<script lang="ts">
    import MySelectCard from "../../../component/card/MySelectCard.svelte";
    import {
        select_homepage,
        current_homepage,
        homepage_list,
        homepage_url,
        homepage_cache,
        current_view,
    } from "../../../store/changeBody";
    import { slide_opacity, slide_up } from "../../../store/functions";
    import { onDestroy, onMount } from "svelte";
    import MyRadioButton from "../../../component/button/MyRadioButton.svelte";
    import MyTextInput from "../../../component/input/MyTextInput.svelte";
    import {
        ReadConfig,
        WriteConfig,
    } from "../../../../wailsjs/go/launcher/ReaderWriter";
    import MyLoadingPickaxe from "../../../component/card/MyLoadingPickaxe.svelte";
    import RedstoneLampOn from "../../../assets/images/Blocks/RedstoneLampOn.png";
    import MyNormalSpan from "../../../component/input/MyNormalSpan.svelte";
    import {
        GenerateTutorialHomePage,
        GetAllHomePage,
    } from "../../../../wailsjs/go/launcher/MainMethod";
    import MyNormalButton from "../../../component/button/MyNormalButton.svelte";
    import { messagebox, showHint } from "../../../store/messagebox";

    export let slide = null;
    export let after_leave = null;
    let isSetting = false;
    let onlineValue = "";
    let isTransitioning = true;
    let f = false;
    const unsubscribe_current_homepage = select_homepage.subscribe((value) => {
        if (!f) {
            setTimeout(() => {
                f = true;
                isTransitioning = true;
            }, 300);
        } else {
            isTransitioning = !isTransitioning;
            homepage_list.set([]);
            homepage_cache.set("");
            current_homepage.set(0);
            loadingHomePage();
        }
    });
    async function selectAnyHomePage(index: number) {
        isSetting = false;
        select_homepage.set(index);
        await WriteConfig(
            "current",
            "Misc",
            "SelectHomePage",
            index.toString(),
        );
    }
    async function onlineTextInput(event: CustomEvent) {
        let text = event.detail.value;
        let reg =
            /^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\*\+,;=.]+$/gi;
        if (!reg.test(text)) {
            return;
        }
        homepage_url.set(text);
        homepage_cache.set("");
        await WriteConfig("current", "Misc", "HomePageURL", text);
    }
    function control_leave() {
        isTransitioning = true;
    }
    async function anyHomePageClick(index: number) {
        current_homepage.set(index);
        homepage_cache.set("");
        await WriteConfig("current", "Misc", "HomePageValue", index.toString());
    }
    async function loadingHomePage() {
        if ($select_homepage == 1) {
            loading_text = "正在加载预设主页中~";
            loading_state = false;
            setTimeout(() => {
                loading_text =
                    "咱们还不支持预设主页喵~你可以加载一下本地主页或者联网更新！~";
                loading_state = true;
            }, 1000);
            // TODO: 预设主页
        } else if ($select_homepage == 2) {
            loading_text = "正在加载本地主页中~";
            loading_state = false;
            if ($homepage_list.length <= 0) {
                let hps = await GetAllHomePage();
                if (!hps.status) {
                    loading_text = "加载本地主页失败，错误信息：" + hps.message;
                    loading_state = true;
                    return;
                }
                homepage_list.set([...hps.data]);
            }
            isSetting = true;
        }
    }
    async function genTutorialHomePage() {
        await GenerateTutorialHomePage();
        showHint("生成教学主页成功喵~");
    }
    onMount(async () => {
        await loadingHomePage();
        if ($select_homepage < 0) {
            let index = Number(
                await ReadConfig("current", "Misc", "SelectHomePage"),
            );
            select_homepage.set(Number(index >= 0 ? index : -1));
        }
        if ($current_homepage < 0) {
            let index = Number(
                await ReadConfig("current", "Misc", "HomePageValue"),
            );
            current_homepage.set(Number(index >= 0 ? index : -1));
        }
        if ($homepage_url == "") {
            let url = await ReadConfig("current", "Misc", "HomePageURL");
            homepage_url.set(url);
        }
        onlineValue = $homepage_url;
    });
    let loading_state = false;
    let loading_text = "正在加载预设主页中~";
    onDestroy(unsubscribe_current_homepage);
</script>

<div class="component-other" in:slide out:slide on:outroend={after_leave}>
    <MySelectCard title="主页设置">
        <div class="homepage-radio">
            <MyRadioButton
                isChecked={$select_homepage === 0}
                on:click={() => selectAnyHomePage(0)}>空白</MyRadioButton
            >
            <MyRadioButton
                isChecked={$select_homepage === 1}
                on:click={() => selectAnyHomePage(1)}>预设</MyRadioButton
            >
            <MyRadioButton
                isChecked={$select_homepage === 2}
                on:click={() => selectAnyHomePage(2)}
                >读取本地文件</MyRadioButton
            >
            <MyRadioButton
                isChecked={$select_homepage === 3}
                on:click={() => selectAnyHomePage(3)}>联网更新</MyRadioButton
            >
        </div>
    </MySelectCard>
    {#if isTransitioning && $select_homepage === 0}
        <div in:slide_up out:slide_up on:outroend={control_leave}></div>
    {:else if isTransitioning && $select_homepage === 1}
        {#if isTransitioning && isSetting}
            <div in:slide_up out:slide_up on:outroend={control_leave}>
                <MySelectCard title="预设主页">
                    <div class="version-all">
                        {#each $homepage_list as phl, i}
                            <div
                                class="a-homepage"
                                on:click={() => {
                                    anyHomePageClick(i);
                                }}
                                on:keydown|preventDefault
                                style={$current_homepage === i
                                    ? "cursor: default"
                                    : "cursor: pointer"}
                            >
                                <MyRadioButton
                                    isChecked={i === $current_homepage}
                                    style_in="margin-left: 5px"
                                />
                                <img
                                    src={RedstoneLampOn}
                                    alt="主页"
                                    class="a-homepage-icon"
                                />
                                <div
                                    class="info"
                                    style="display: flex; flex-direction: column; margin-left: 10px; pointer-events: none"
                                >
                                    <MyNormalSpan>{phl.name}</MyNormalSpan>
                                    <MyNormalSpan
                                        style_in="font-size: 14px; color: gray;"
                                        >{phl.path}</MyNormalSpan
                                    >
                                </div>
                            </div>
                        {/each}
                    </div>
                </MySelectCard>
            </div>
        {:else if isTransitioning}
            <div
                in:slide_opacity
                out:slide_opacity
                on:outroend={control_leave}
                style="width: 100%; height: calc(100% - 120px); display: flex; flex-direction: column; align-items: center; justify-content: center;"
            >
                <MyLoadingPickaxe {loading_text} state={loading_state} />
            </div>
        {/if}
    {:else if isTransitioning && $select_homepage === 2}
        {#if isTransitioning && isSetting}
            <div in:slide_up out:slide_up on:outroend={control_leave}>
                <MySelectCard title="本地主页">
                    <div class="version-all">
                        <div
                            class="hint blue-hint"
                            style="width: calc(100% - 30px);"
                        >
                            从 PCL.Nova/homepage/&lt;文件名&gt;.nxml
                            读取主页内容。<br />
                            你可以手动编辑该文件，向主页添加文本、图片、常用网站、快捷启动等功能。
                        </div>
                        <div
                            class="yellow-hint hint"
                            style="width: calc(100% - 30px);"
                        >
                            请谨慎使用陌生人提供的主页，恶意代码可能会造成损害！
                        </div>
                        <div style="margin-top: 5px; margin-bottom: 5px;">
                            <MyNormalButton
                                style_in="width: 100px; height: 40px;"
                                on:click={() => {
                                    homepage_cache.set("");
                                    current_view.set("home");
                                }}>刷新主页</MyNormalButton
                            >
                            <MyNormalButton
                                style_in="width: 100px; height: 40px; margin-left: 20px;"
                                on:click={genTutorialHomePage}
                                >生成教学主页</MyNormalButton
                            >
                            <MyNormalButton
                                style_in="width: 100px; height: 40px; margin-left: 20px;"
                                on:click={() => {
                                    (async function xx() {
                                        await messagebox(
                                            "主页自定义教程",
                                            "1. 点击 生成教学主页 按钮，这会在 PCL.Nova/homepage/ 文件夹下生成 Custom.nxml 布局文件。<br>2. 使用记事本等工具打开这个文件并进行修改，修改完记得保存。<br>3. 点击刷新主页按钮，查看主页现在长啥样了。<br><br>你可以在生成教学文件后直接刷新主页，对照着进行修改，更有助于理解。",
                                        );
                                    })();
                                }}>查看教程</MyNormalButton
                            >
                        </div>
                        {#each $homepage_list as phl, i}
                            <div
                                class="a-homepage"
                                on:click={() => {
                                    anyHomePageClick(i);
                                }}
                                on:keydown|preventDefault
                                style={$current_homepage === i
                                    ? "cursor: default"
                                    : "cursor: pointer"}
                            >
                                <MyRadioButton
                                    isChecked={i === $current_homepage}
                                    style_in="margin-left: 5px"
                                />
                                <img
                                    src={RedstoneLampOn}
                                    alt="主页"
                                    class="a-homepage-icon"
                                />
                                <div
                                    class="info"
                                    style="display: flex; flex-direction: column; margin-left: 10px; pointer-events: none"
                                >
                                    <MyNormalSpan>{phl.name}</MyNormalSpan>
                                    <MyNormalSpan
                                        style_in="font-size: 14px; color: gray;"
                                        >{phl.path}</MyNormalSpan
                                    >
                                </div>
                            </div>
                        {/each}
                    </div>
                </MySelectCard>
            </div>
        {:else if isTransitioning}
            <div
                in:slide_opacity
                out:slide_opacity
                on:outroend={control_leave}
                style="width: 100%; height: calc(100% - 120px); display: flex; flex-direction: column; align-items: center; justify-content: center;"
            >
                <MyLoadingPickaxe {loading_text} state={loading_state} />
            </div>
        {/if}
    {:else if isTransitioning && $select_homepage === 3}
        <div in:slide_up out:slide_up on:outroend={control_leave}>
            <MySelectCard title="联网主页">
                <div class="version-all">
                    <div
                        class="hint blue-hint"
                        style="width: calc(100% - 30px);"
                    >
                        从指定网址联网获取主页内容，服主也可以用于动态更新服务器公告。
                    </div>
                    <div
                        class="yellow-hint hint"
                        style="width: calc(100% - 30px)"
                    >
                        请谨慎使用陌生人提供的主页，恶意代码可能会造成损害！
                    </div>
                    <div
                        style="display: flex; align-items: center; margin-top: 5px;"
                    >
                        <div>下载地址</div>
                        <div style="flex: 1; margin-left: 50px">
                            <MyTextInput
                                style_in="height: 30px; width: calc(100% - 30px);"
                                bind:value={onlineValue}
                                on:blur={onlineTextInput}
                            />
                        </div>
                    </div>
                </div>
            </MySelectCard>
        </div>
    {/if}
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
    .homepage-radio {
        display: flex;
        align-items: center;
        justify-content: space-around;
        margin: 10px 0;
    }
    .a-homepage {
        position: relative;
        transition: all 0.2s;
        border-radius: 10px;
        height: 50px;
        width: calc(100% - 2px);
        flex-shrink: 0;
        display: flex;
        align-items: center;
    }
    .a-homepage-icon {
        width: 40px;
        height: 40px;
        image-rendering: pixelated;
        border-radius: 10px;
        box-shadow: 0 0 5px gray;
        margin-left: 5px;
    }
    .a-homepage:hover {
        background-color: rgba(128, 128, 128, 0.5);
    }
    .a-homepage:active {
        transform: scale(99%);
    }
</style>
