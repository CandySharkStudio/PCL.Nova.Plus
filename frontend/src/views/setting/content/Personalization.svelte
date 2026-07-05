<script lang="ts">
    import MyToggleSwitch from "../../../component/button/MyToggleSwitch.svelte";

    export let slide = null;
    export let after_leave = null;
    import { theme_mode, dark_mode } from "../../../store/changeBody";
    import MyNormalButton from "../../../component/button/MyNormalButton.svelte";
    import MyNormalSpan from "../../../component/input/MyNormalSpan.svelte";
    import MySelectCard from "../../../component/card/MySelectCard.svelte";
    import MyRadioButton from "../../../component/button/MyRadioButton.svelte";
    import { OpenCustomURL } from "../../../store/functions";
    import {
        WriteConfig,
        ReadConfig,
    } from "../../../../wailsjs/go/launcher/ReaderWriter";
    import { unlock_theme } from "../../../store/changeBody";
    async function changeTheme(index: number) {
        theme_mode.set(index);
        await WriteConfig(
            "other",
            "Unlock",
            "ThemeMode",
            $theme_mode.toString(),
        );
    }
    async function toggleDark(is: boolean) {
        dark_mode.set(!is);
        await WriteConfig(
            "current",
            "Misc",
            "DarkMode",
            $dark_mode ? "1" : "0",
        );
    }
</script>

<div class="component-person" in:slide out:slide on:outroend={after_leave}>
    <MySelectCard title="主题" isExpand={false}>
        <div class="proc" style="position: relative;">
            {#if !$unlock_theme}
                <div id="mask">
                    <MyNormalSpan style_in="font-weight: bold"
                        >请支持官方版本以使用主题功能</MyNormalSpan
                    >
                    <MyNormalButton
                        style_in="width: 200px; height: 40px"
                        on:click={() =>
                            OpenCustomURL("https://afdian.com/a/LTCat")}
                        >点我进入龙猫的爱发电</MyNormalButton
                    >
                </div>
            {/if}
            <div class="line">
                <MyRadioButton
                    isChecked={$theme_mode === 1}
                    on:click={() => changeTheme(1)}>龙猫蓝</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 2}
                    on:click={() => changeTheme(2)}>甜柠青</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 3}
                    on:click={() => changeTheme(3)}>小草绿</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 4}
                    on:click={() => changeTheme(4)}>菠萝黄</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 5}
                    on:click={() => changeTheme(5)}>橡木棕</MyRadioButton
                >
            </div>
            <div class="line">
                <MyRadioButton
                    isChecked={$theme_mode === 6}
                    on:click={() => changeTheme(6)}>玄素黑</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 7}
                    on:click={() => changeTheme(7)}>滑稽彩</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 8}
                    on:click={() => changeTheme(8)}>铁杆粉</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 9}
                    on:click={() => changeTheme(9)}>神秘紫</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 10}
                    on:click={() => changeTheme(10)}>欧皇彩</MyRadioButton
                >
            </div>
            <div class="line">
                <MyRadioButton
                    isChecked={$theme_mode === 11}
                    on:click={() => changeTheme(11)}>秋仪金</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 12}
                    on:click={() => changeTheme(12)}>活跃橙</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 13}
                    on:click={() => changeTheme(13)}>跳票红</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 14}
                    on:click={() => changeTheme(14)}>极客蓝</MyRadioButton
                >
                <MyRadioButton
                    isChecked={$theme_mode === 15}
                    on:click={() => changeTheme(15)}>自定义</MyRadioButton
                >
            </div>
        </div>
    </MySelectCard>
    <MySelectCard title="个性化" isExpand={true} canExpand={true}>
        <div class="proc">
            <div id="dark">
                <MyNormalSpan style_in="font-weight: bold"
                    >暗色模式</MyNormalSpan
                >
                <MyToggleSwitch
                    isSelect={$dark_mode}
                    on:click={() => toggleDark($dark_mode)}
                />
            </div>
        </div>
    </MySelectCard>
</div>

<style>
    .component-person {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        overflow-y: auto;
    }
    .line {
        width: 100%;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: space-around;
    }
    #mask {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        backdrop-filter: blur(2px);
        z-index: 2;
    }
    #dark {
        width: 100%;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
</style>
