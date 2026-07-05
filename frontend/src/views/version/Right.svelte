<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import Grass from "../../assets/images/Blocks/Grass.png";
    export let slide = null;
    export let after_leave = null;
    import {
        select_mc,
        select_mc_version,
        current_mc_index,
        current_mc_version_index,
        current_mc_version_path,
    } from "../../store/mc";
    import {
        GetCurrentMinecraftDir,
        GetMCAllVersion,
        GetMCVersionConfig,
        SetMCVersionConfig,
    } from "../../../wailsjs/go/launcher/LaunchMethod";
    import MySelectCard from "../../component/card/MySelectCard.svelte";
    import MyLoadingPickaxe from "../../component/card/MyLoadingPickaxe.svelte";
    import MyCardButton from "../../component/button/MyCardButton.svelte";
    import { current_view } from "../../store/changeBody";
    import MyNormalButton from "../../component/button/MyNormalButton.svelte";
    import { WriteConfig } from "../../../wailsjs/go/launcher/ReaderWriter";
    import { MSG_WARNING, showHint } from "../../store/messagebox";
    import { openExplorer } from "../../store/functions";
    import { inputbox, messagebox } from "../../store/messagebox.js";
    import { launcher } from "../../../wailsjs/go/models";
    import MyNormalSpan from "../../component/input/MyNormalSpan.svelte";
    let fpath = "";
    let isTransitioning = true;
    let f = false;
    let show = false;
    let height = 0;
    async function selectMCVersion(num: number) {
        let c_select_mc_length = [];
        if ($select_mc.length <= 0) {
            let v = await GetMCVersionConfig();
            if (!v.status) {
                showHint(
                    "解析 MC 版本路径时出现了点问题，错误信息：" + v.message,
                );
                return;
            }
            for (let i = 0; i < v.data.mc.length; i++) {
                c_select_mc_length.push({
                    path: v.data.mc[i].path,
                    name: v.data.mc[i].name,
                });
            }
        } else {
            c_select_mc_length = $select_mc;
        }
        let rp =
            num <= 0
                ? await GetCurrentMinecraftDir()
                : c_select_mc_length[num - 1].path;
        fpath = rp;
        let paths = await GetMCAllVersion(rp);
        height = 0;
        if (paths.length == 0) {
            show = true;
        } else {
            show = false;
            select_mc_version.set([]);
            for (let i = 0; i < paths.length; i++) {
                let p1 = paths[i].lastIndexOf("\\");
                let p2 = paths[i].lastIndexOf("/");
                let p3 =
                    p1 > 1
                        ? paths[i].substring(p1 + 1)
                        : p2 > 0
                          ? paths[i].substring(p2 + 1)
                          : "";
                height += 50;
                select_mc_version.set([
                    ...$select_mc_version,
                    {
                        path: paths[i],
                        name: p3,
                    },
                ]);
            }
        }
        height += 20;
    }
    async function selectOneVersion(num: number) {
        await WriteConfig("current", "MC", "SelectVer", num.toString());
        current_mc_version_index.set(num);
        current_mc_version_path.set($select_mc_version[num].path);
        current_view.set("home");
    }
    async function openMCInstanceSetting(num: number) {
        current_mc_version_index.set(num);
        current_view.set("instance");
    }
    function control_leave() {
        isTransitioning = true;
    }
    const unsubscribe_current_mc_index = current_mc_index.subscribe((value) => {
        if (!f) {
            f = true;
            isTransitioning = true;
        } else {
            isTransitioning = !isTransitioning;
        }
        selectMCVersion(value);
    });
    onDestroy(unsubscribe_current_mc_index);
    onMount(() => {
        selectMCVersion($current_mc_index);
    });
    async function renameName() {
        let ren = await inputbox("重命名名称", "请输入要重命名的名称~");
        if (ren == "") return;
        select_mc.update((value) => {
            value[$current_mc_index - 1] = { path: fpath, name: ren };
            return value;
        });
        await SetMCVersionConfig(
            launcher.MCConfigs.createFrom({ mc: $select_mc }),
        );
    }
    async function removeName() {
        if (
            (await messagebox(
                "是否确定移除",
                "你真的确定要移除这个版本文件夹吗？此时只是相对于 PCL2.Nova 移除了，但是文件资源管理器里面绝对还在！所以不用担心~",
                MSG_WARNING,
                ["yes", "no"],
            )) == 1
        ) {
            return;
        }
        select_mc.update((value) => {
            value.splice($current_mc_index - 1, 1);
            return value;
        });
        current_mc_index.set(0);
        await WriteConfig(
            "current",
            "MC",
            "SelectMC",
            $current_mc_index.toString(),
        );
        await SetMCVersionConfig(
            launcher.MCConfigs.createFrom({ mc: $select_mc }),
        );
    }
</script>

<div class="component-right" in:slide out:slide on:outroend={after_leave}>
    {#if isTransitioning}
        <div id="transitioning" in:slide out:slide on:outroend={control_leave}>
            <div id="version-name">
                <MySelectCard
                    title="当前文件夹"
                    isExpand={true}
                    in_style="margin-bottom: 20px"
                >
                    <div class="version-all">
                        <MyNormalSpan>{fpath}</MyNormalSpan><br />
                        <MyNormalButton
                            style_in="width: 100px; height: 30px; margin-right: 10px"
                            on:click={() => openExplorer(fpath)}
                        >
                            打开文件夹
                        </MyNormalButton>
                        {#if $current_mc_index !== 0}
                            <MyNormalButton
                                style_in="width: 100px; height: 30px; margin-right: 10px"
                                on:click={renameName}
                            >
                                重命名文件夹
                            </MyNormalButton>
                            <MyNormalButton
                                style_in="width: 100px; height: 30px; margin-right: 10px; border: 1px solid red; color: red;"
                                on:click={removeName}
                            >
                                移除文件夹
                            </MyNormalButton>
                        {/if}
                    </div>
                </MySelectCard>
                {#if show}
                    <div class="loading">
                        <MyLoadingPickaxe
                            state={true}
                            loading_text="未找到任何版本，请下载一个先吧~"
                        />
                    </div>
                {:else}
                    <div class="scroll-div">
                        <MySelectCard
                            title="MC 实例"
                            in_style="margin-top: 1px"
                            maxHeight={height}
                        >
                            <div class="version-all">
                                {#each $select_mc_version as ver, i}
                                    <MyCardButton
                                        image={Grass}
                                        title={ver.name}
                                        desc={ver.name}
                                        click={() => selectOneVersion(i)}
                                        buttons={[
                                            {
                                                title: "设置",
                                                click: () =>
                                                    // openMCInstanceSetting(i),
                                                    {},
                                                icon: `
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                viewBox="0 0 24 24"
                                                fill="none"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="1.5"
                                                color="currentColor">
                                                <path d="m21.318 7.141l-.494-.856c-.373-.648-.56-.972-.878-1.101c-.317-.13-.676-.027-1.395.176l-1.22.344c-.459.106-.94.046-1.358-.17l-.337-.194a2 2 0 0 1-.788-.967l-.334-.998c-.22-.66-.33-.99-.591-1.178c-.261-.19-.609-.19-1.303-.19h-1.115c-.694 0-1.041 0-1.303.19c-.261.188-.37.518-.59 1.178l-.334.998a2 2 0 0 1-.789.967l-.337.195c-.418.215-.9.275-1.358.17l-1.22-.345c-.719-.203-1.078-.305-1.395-.176c-.318.129-.505.453-.878 1.1l-.493.857c-.35.608-.525.911-.491 1.234c.034.324.268.584.736 1.105l1.031 1.153c.252.319.431.875.431 1.375s-.179 1.056-.43 1.375l-1.032 1.152c-.468.521-.702.782-.736 1.105s.14.627.49 1.234l.494.857c.373.647.56.971.878 1.1s.676.028 1.395-.176l1.22-.344a2 2 0 0 1 1.359.17l.336.194c.36.23.636.57.788.968l.334.997c.22.66.33.99.591 1.18c.262.188.609.188 1.303.188h1.115c.694 0 1.042 0 1.303-.189s.371-.519.59-1.179l.335-.997c.152-.399.428-.738.788-.968l.336-.194c.42-.215.9-.276 1.36-.17l1.22.344c.718.204 1.077.306 1.394.177c.318-.13.505-.454.878-1.101l.493-.857c.35-.607.525-.91.491-1.234s-.268-.584-.736-1.105l-1.031-1.152c-.252-.32-.431-.875-.431-1.375s.179-1.056.43-1.375l1.032-1.153c.468-.52.702-.781.736-1.105s-.14-.626-.49-1.234"/>
                                                <path d="M15.52 12a3.5 3.5 0 1 1-7 0a3.5 3.5 0 0 1 7 0"/>
                                            </svg>
                                            `,
                                            },
                                        ]}
                                    />
                                {/each}
                            </div>
                        </MySelectCard>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .component-right {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }
    #transitioning {
        width: 100%;
        height: 100%;
    }
    #version-name {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
    }
    .loading {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .scroll-div {
        flex: 1;
        overflow-y: auto;
    }
</style>
