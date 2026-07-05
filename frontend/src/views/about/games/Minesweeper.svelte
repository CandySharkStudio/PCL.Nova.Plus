<script lang="ts">
    import MyTextInput from "../../../component/input/MyTextInput.svelte";
    import MyNormalSpan from "../../../component/input/MyNormalSpan.svelte";
    import MyNormalButton from "../../../component/button/MyNormalButton.svelte";
    import MyToggleSwitch from "../../../component/button/MyToggleSwitch.svelte";
    export let slide = null;
    export let after_leave = null;

    // 格子接口
    interface Grids {
        // 格子横坐标
        x: number;
        // 格子纵坐标
        y: number;
        // 格子数字（-1：雷、0：周围8格没有雷、1：有一个雷（剩下以此类推）、2、3、4、5、6、7、8）
        p: -1 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
        // 格子是否被标记（-1：已被点击、0：未标记、1：旗子、2：问号）
        // 3比较特殊，用于在开启扣分模式后，如果点到了雷则默认将该旗子标记为3，用于表示该格子不可点击，但是却被标了旗子。
        m: -1 | 0 | 1 | 2 | 3;
        // 格子应该被显示成什么（有炸弹、旗子、问号、1~8的数字，其中0不会显示）
        s: string;
        // 是否是当前点击的雷（用于判断点击该雷的时候背景变成红色。）
        c: boolean;
    }

    let score = 0;
    let time = 0;
    let flags = 0;
    let win = 0;

    let width = "";
    let height = "";
    let mines = "";

    // 格子列表
    let grids: Grids[][] = [];

    // 临时宽度
    let tempWidth = 0;
    // 临时高度
    let tempHeight = 0;
    // 临时雷数
    let tempMines = 0;
    // 棋盘记录雷数
    let mine_count = 0;

    // 是否开始游戏
    let start = false;
    // 是否锁住棋盘（无法被点击）
    let locked = false;
    // 扣分模式
    let cheat = false;

    setInterval(() => {
        if (start && !locked) {
            time++;
        }
    }, 1000);
    //初始化数组
    function init_array() {
        grids = [];
        for (let i = 0; i < tempHeight; i++) {
            grids[i] = [];
            for (let j = 0; j < tempWidth; j++) {
                grids[i][j] = {
                    x: j,
                    y: i,
                    p: 0,
                    m: 0,
                    s: "",
                    c: false,
                };
            }
        }
    }
    //生成雷
    function generate_mine() {
        for (let i = 0; i < tempMines; i++) {
            while (true) {
                let x = Math.floor(Math.random() * tempHeight);
                let y = Math.floor(Math.random() * tempWidth);
                if (grids[x][y].p == -1) {
                    continue;
                }
                grids[x][y].p = -1;
                break;
            }
        }
    }
    //初始化数字（为雷的周围生成数字）
    function init_number() {
        for (let i = 0; i < tempHeight; i++) {
            for (let j = 0; j < tempWidth; j++) {
                if (grids[i][j].p == -1) {
                    continue;
                }
                let foo = 0;
                if (i > 0 && j > 0) {
                    if (grids[i - 1][j - 1].p == -1) {
                        foo++;
                    }
                }
                if (i > 0 && j < tempWidth - 1) {
                    if (grids[i - 1][j + 1].p == -1) {
                        foo++;
                    }
                }
                if (i < tempHeight - 1 && j > 0) {
                    if (grids[i + 1][j - 1].p == -1) {
                        foo++;
                    }
                }
                if (i < tempHeight - 1 && j < tempWidth - 1) {
                    if (grids[i + 1][j + 1].p == -1) {
                        foo++;
                    }
                }
                if (i > 0) {
                    if (grids[i - 1][j].p == -1) {
                        foo++;
                    }
                }
                if (i < tempHeight - 1) {
                    if (grids[i + 1][j].p == -1) {
                        foo++;
                    }
                }
                if (j > 0) {
                    if (grids[i][j - 1].p == -1) {
                        foo++;
                    }
                }
                if (j < tempWidth - 1) {
                    if (grids[i][j + 1].p == -1) {
                        foo++;
                    }
                }
                grids[i][j].p = foo as 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
            }
        }
    }
    function startGame() {
        let w = parseInt(width);
        let h = parseInt(height);
        let m = parseInt(mines);
        if (
            Number.isNaN(w) ||
            Number.isNaN(h) ||
            Number.isNaN(m) ||
            w < 5 ||
            w > 100 ||
            h < 5 ||
            h > 100 ||
            m < 5 ||
            m > w * h - 10
        ) {
            return;
        }
        tempWidth = w;
        tempHeight = h;
        tempMines = m;
        init_array();
        generate_mine();
        init_number();
        start = false;
        time = 0;
        flags = m;
        win = 0;
        score = 0;
        locked = false;
        mine_count = 0;
    }

    //格子左键点击（揭开旗子）
    function grid_button_click(x: number, y: number) {
        if (locked) return;
        if (!start) start = true;
        if (![-1, 1, 3].includes(grids[y][x].m)) {
            grids[y][x].m = -1;
            grids[y][x].s =
                grids[y][x].p == -1
                    ? `<svg xmlns="http://www.w3.org/2000/svg" style="width: 100%; height: 100%;" viewBox="0 0 32 32"><g fill="none"><path fill="#FF822D" d="M23.369 2.555c.18.728 1.208.745 1.413.024a.73.73 0 0 1 .994-.47l.585.254a.73.73 0 0 1 .335 1.048c-.386.642.328 1.38.982 1.016a.73.73 0 0 1 1.036.37l.235.594a.73.73 0 0 1-.504.978c-.727.18-.745 1.208-.024 1.413a.73.73 0 0 1 .47.994l-.254.585a.73.73 0 0 1-1.048.335c-.642-.386-1.38.328-1.016.982a.73.73 0 0 1-.37 1.036l-.594.235a.73.73 0 0 1-.978-.504c-.18-.728-1.208-.745-1.413-.024a.73.73 0 0 1-.994.47l-.585-.254a.73.73 0 0 1-.335-1.048c.386-.642-.328-1.38-.982-1.016a.73.73 0 0 1-1.036-.37l-.235-.594a.73.73 0 0 1 .504-.978c.728-.18.745-1.208.024-1.413a.73.73 0 0 1-.47-.994l.254-.585a.73.73 0 0 1 1.048-.335c.642.386 1.38-.328 1.016-.983a.73.73 0 0 1 .37-1.035l.594-.235a.73.73 0 0 1 .978.504"/><path fill="#F4F4F4" d="M25.207 5.793a1 1 0 0 1 0 1.414l-3 3a1 1 0 0 1-1.414-1.414l3-3a1 1 0 0 1 1.414 0"/><path fill="#533566" d="M26 18c0 6.627-5.373 12-12 12S2 24.627 2 18S7.373 6 14 6s12 5.373 12 12"/><path fill="#6B438B" d="M12 27c6.075 0 11-4.925 11-11a11 11 0 0 0-.489-3.252A9.95 9.95 0 0 1 24 18c0 5.523-4.477 10-10 10a9.95 9.95 0 0 1-5.252-1.489C9.776 26.83 10.868 27 12 27"/></g></svg>`
                    : grids[y][x].p == 0
                      ? ""
                      : grids[y][x].p.toString();
            if (grids[y][x].p == -1) {
                if (!cheat) {
                    locked = true;
                    win = 2;
                    for (let k = 0; k < tempHeight; k++) {
                        for (let l = 0; l < tempWidth; l++) {
                            if (grids[k][l].p == -1) {
                                grids[k][l].s = `<svg xmlns="http://www.w3.org/2000/svg" style="width: 100%; height: 100%;" viewBox="0 0 32 32"><g fill="none"><path fill="#FF822D" d="M23.369 2.555c.18.728 1.208.745 1.413.024a.73.73 0 0 1 .994-.47l.585.254a.73.73 0 0 1 .335 1.048c-.386.642.328 1.38.982 1.016a.73.73 0 0 1 1.036.37l.235.594a.73.73 0 0 1-.504.978c-.727.18-.745 1.208-.024 1.413a.73.73 0 0 1 .47.994l-.254.585a.73.73 0 0 1-1.048.335c-.642-.386-1.38.328-1.016.982a.73.73 0 0 1-.37 1.036l-.594.235a.73.73 0 0 1-.978-.504c-.18-.728-1.208-.745-1.413-.024a.73.73 0 0 1-.994.47l-.585-.254a.73.73 0 0 1-.335-1.048c.386-.642-.328-1.38-.982-1.016a.73.73 0 0 1-1.036-.37l-.235-.594a.73.73 0 0 1 .504-.978c.728-.18.745-1.208.024-1.413a.73.73 0 0 1-.47-.994l.254-.585a.73.73 0 0 1 1.048-.335c.642.386 1.38-.328 1.016-.983a.73.73 0 0 1 .37-1.035l.594-.235a.73.73 0 0 1 .978.504"/><path fill="#F4F4F4" d="M25.207 5.793a1 1 0 0 1 0 1.414l-3 3a1 1 0 0 1-1.414-1.414l3-3a1 1 0 0 1 1.414 0"/><path fill="#533566" d="M26 18c0 6.627-5.373 12-12 12S2 24.627 2 18S7.373 6 14 6s12 5.373 12 12"/><path fill="#6B438B" d="M12 27c6.075 0 11-4.925 11-11a11 11 0 0 0-.489-3.252A9.95 9.95 0 0 1 24 18c0 5.523-4.477 10-10 10a9.95 9.95 0 0 1-5.252-1.489C9.776 26.83 10.868 27 12 27"/></g></svg>`;
                            }
                        }
                    }
                    grids[y][x].c = true;
                    return;
                } else {
                    grids[y][x].c = true;
                    score -= 50;
                    grids[y][x].m = 3;
                    flags -= 1;
                }
            } else {
                mine_count++;
                score += 10;
            }
            if (mine_count == tempHeight * tempWidth - tempMines) {
                locked = true;
                win = 1;
                return;
            }
            if (grids[y][x].p != 0) return;
            if (x > 0 && y > 0) {
                grid_button_click(x - 1, y - 1);
            }
            if (x > 0 && y < tempHeight - 1) {
                grid_button_click(x - 1, y + 1);
            }
            if (x < tempWidth - 1 && y > 0) {
                grid_button_click(x + 1, y - 1);
            }
            if (x < tempWidth - 1 && y < tempHeight - 1) {
                grid_button_click(x + 1, y + 1);
            }
            if (x > 0) {
                grid_button_click(x - 1, y);
            }
            if (x < tempWidth - 1) {
                grid_button_click(x + 1, y);
            }
            if (y > 0) {
                grid_button_click(x, y - 1);
            }
            if (y < tempHeight - 1) {
                grid_button_click(x, y + 1);
            }
        }
    }

    //格子右键点击（标旗或者揭开【周围】格子）
    function grid_button_right(x: number, y: number) {
        if (locked) return;
        switch (grids[y][x].m) {
            case -1:
                let cc = grids[y][x].p;
                let k = 0;
                if (x > 0 && y > 0)
                    if ([1, 3].includes(grids[y - 1][x - 1].m)) k++;
                if (x > 0 && y < tempHeight - 1)
                    if ([1, 3].includes(grids[y + 1][x - 1].m)) k++;
                if (x < tempWidth - 1 && y > 0)
                    if ([1, 3].includes(grids[y - 1][x + 1].m)) k++;
                if (x < tempWidth - 1 && y < tempHeight - 1)
                    if ([1, 3].includes(grids[y + 1][x + 1].m)) k++;
                if (x > 0) if ([1, 3].includes(grids[y][x - 1].m)) k++;
                if (x < tempWidth - 1)
                    if ([1, 3].includes(grids[y][x + 1].m)) k++;
                if (y > 0) if ([1, 3].includes(grids[y - 1][x].m)) k++;
                if (y < tempHeight - 1)
                    if ([1, 3].includes(grids[y + 1][x].m)) k++;
                // 如果格子本身的数字小于或者等于周围插旗子的数字，则开启
                if (cc > k) return;
                if (x > 0 && y > 0) grid_button_click(x - 1, y - 1);
                if (x > 0 && y < tempHeight - 1)
                    grid_button_click(x - 1, y + 1);
                if (x < tempWidth - 1 && y > 0) grid_button_click(x + 1, y - 1);
                if (x < tempWidth - 1 && y < tempHeight - 1)
                    grid_button_click(x + 1, y + 1);
                if (x > 0) grid_button_click(x - 1, y);
                if (x < tempWidth - 1) grid_button_click(x + 1, y);
                if (y > 0) grid_button_click(x, y - 1);
                if (y < tempHeight - 1) grid_button_click(x, y + 1);
                break;
            case 0:
                if (flags > 0) {
                    grids[y][x].m = 1;
                    grids[y][x].s = `<svg xmlns="http://www.w3.org/2000/svg" style="width: 100%; height: 100%;" viewBox="0 0 32 32"><g fill="none"><path fill="#F8312F" d="M27.626 11.944L8 4H7v17h1l19.626-7.944a.6.6 0 0 0 0-1.112"/><path fill="#E39D89" d="M4 4a2 2 0 1 1 4 0v26H4z"/></g></svg>`;
                    flags--;
                } else {
                    grids[y][x].m = 2;
                    grids[y][x].s = `<svg xmlns="http://www.w3.org/2000/svg" style="width: 100%; height: 100%;" viewBox="0 0 32 32"><path fill="#9B9B9B" d="M13.037 20.863c0 1.302 1.145 2.36 2.555 2.36s2.556-1.058 2.556-2.37V19.15c0-.111.073-.209.18-.242C21.676 17.85 24 14.919 24 11.562v-1.254c0-4.239-3.69-7.773-8.227-7.861c-2.28-.05-4.432.744-6.065 2.212c-1.622 1.469-2.523 3.447-2.523 5.552c0 1.311 1.155 2.369 2.566 2.369s2.555-1.058 2.555-2.36c0-.822.35-1.596.986-2.173a3.4 3.4 0 0 1 2.375-.872c1.78.04 3.223 1.45 3.223 3.143v1.244c0 1.468-1.124 2.731-2.683 2.996c-1.834.313-3.17 1.791-3.17 3.514zM15.5 30a2.5 2.5 0 1 0 0-5a2.5 2.5 0 0 0 0 5"/></svg>`;
                }
                break;
            case 1:
                grids[y][x].m = 2;
                grids[y][x].s = `<svg xmlns="http://www.w3.org/2000/svg" style="width: 100%; height: 100%;" viewBox="0 0 32 32"><path fill="#9B9B9B" d="M13.037 20.863c0 1.302 1.145 2.36 2.555 2.36s2.556-1.058 2.556-2.37V19.15c0-.111.073-.209.18-.242C21.676 17.85 24 14.919 24 11.562v-1.254c0-4.239-3.69-7.773-8.227-7.861c-2.28-.05-4.432.744-6.065 2.212c-1.622 1.469-2.523 3.447-2.523 5.552c0 1.311 1.155 2.369 2.566 2.369s2.555-1.058 2.555-2.36c0-.822.35-1.596.986-2.173a3.4 3.4 0 0 1 2.375-.872c1.78.04 3.223 1.45 3.223 3.143v1.244c0 1.468-1.124 2.731-2.683 2.996c-1.834.313-3.17 1.791-3.17 3.514zM15.5 30a2.5 2.5 0 1 0 0-5a2.5 2.5 0 0 0 0 5"/></svg>`;
                flags++;
                break;
            case 2:
                grids[y][x].m = 0;
                grids[y][x].s = "";
                break;
            default:
                break;
        }
    }
</script>

<div class="component-minesweeper" in:slide out:slide on:outroend={after_leave}>
    <div class="bar">
        <div style:width="100%">
            <MyNormalSpan>场地宽度</MyNormalSpan>
        </div>
        <MyTextInput
            placeholder="[5, 100] 区间范围"
            value={width}
            on:blur={(e) => (width = e.detail.value)}
            style_in="width: 180px; height: 25px; font-size: 15px; transition: all 0.2s;"
        />
        <div style:width="100%">
            <MyNormalSpan>场地高度</MyNormalSpan>
        </div>
        <MyTextInput
            placeholder="[5, 100] 区间范围"
            value={height}
            on:blur={(e) => (height = e.detail.value)}
            style_in="width: 180px; height: 25px; font-size: 15px; transition: all 0.2s;"
        />
        <div style:width="100%">
            <MyNormalSpan>场地雷数</MyNormalSpan>
        </div>
        <MyTextInput
            placeholder="[5, w*h-10] 区间范围"
            value={mines}
            on:blur={(e) => (mines = e.detail.value)}
            style_in="width: 180px; height: 25px; font-size: 15px; transition: all 0.2s;"
        />
        <div style="width: 100%; height: 30px; display: flex; align-items: center; justify-content: center; gap: 5px; margin-top: 5px;">
            <MyNormalButton
                style_in="flex: 1; height: 25px;"
                on:click={() => {
                    width = "9";
                    height = "9";
                    mines = "10";
                }}
            >
                默认条件
            </MyNormalButton>
            <MyNormalButton
                style_in="flex: 1; height: 25px;"
                on:click={startGame}
            >
                开始游戏
            </MyNormalButton>
        </div>
        <div
            style="width: 100%; display: flex; align-items: center; height: 30px; justify-content: space-around; margin-top: 10px"
        >
            <MyNormalSpan>扣分模式</MyNormalSpan>
            <MyToggleSwitch
                on:click={() => (cheat = !cheat)}
                isSelect={cheat}
            />
        </div>
        <div class="info">
            <MyNormalSpan style_in="margin-left: 10px; font-size: 24px;"
                >分数：{score}</MyNormalSpan
            >
            <MyNormalSpan style_in="margin-left: 10px; font-size: 24px;"
                >时间：{time}</MyNormalSpan
            >
            <MyNormalSpan style_in="margin-left: 10px; font-size: 24px;"
                >旗子：{flags}</MyNormalSpan
            >
            <MyNormalSpan
                style_in={"font-size: 50px; " +
                    (win == 1
                        ? "color: green;"
                        : win == 2
                          ? "color: red;"
                          : "")}
                >{win == 1 ? "胜利" : win == 2 ? "失败" : ""}</MyNormalSpan
            >
        </div>
    </div>
    <div class="chess">
        {#each grids as grid1}
            {#each grid1 as grid}
                <button
                    class={[-1, 3].includes(grid.m)
                        ? grid.c
                            ? "grid-button-red"
                            : "grid-button-click"
                        : "grid-button"}
                    style={"left: " +
                        grid.x * 25 +
                        "px; top: " +
                        grid.y * 25 +
                        "px; color: " +
                        (grid.p == 1
                            ? "rgb(0, 128, 255)"
                            : grid.p == 2
                              ? "rgb(0, 128, 0)"
                              : grid.p == 3
                                ? "rgb(192, 0, 0)"
                                : grid.p == 4
                                  ? "rgb(0, 0, 96)"
                                  : grid.p == 5
                                    ? "rgb(128, 0, 48)"
                                    : grid.p == 6
                                      ? "rgb(0, 96, 96)"
                                      : grid.p == 7
                                        ? "rgb(10, 10, 10)"
                                        : grid.p == 8
                                          ? "rgb(100, 100, 100)"
                                          : "black") +
                        ";"}
                    on:click={() => grid_button_click(grid.x, grid.y)}
                    on:contextmenu={() => grid_button_right(grid.x, grid.y)}
                >
                    {@html grid.s}
                </button>
            {/each}
        {/each}
    </div>
</div>

<style>
    .component-minesweeper {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        display: flex;
        overflow: hidden;
    }
    .bar {
        width: 200px;
        height: calc(100% - 20px);
        padding: 10px;
        border-right: 2px solid gray;
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    .info {
        margin-top: 10px;
        display: flex;
        flex-direction: column;
        align-items: start;
        width: 100%;
    }
    .chess {
        width: calc(100% - 222px);
        height: calc(100% - 1px);
        position: relative;
        overflow: auto;
        flex-shrink: 0;
    }
    .grid-button,
    .grid-button-click,
    .grid-button-red {
        position: absolute;
        padding: 2px;
        width: 25px;
        height: 25px;
        border: 1px solid black;
        box-sizing: border-box;
        background-color: lightgray;
        font-size: 16px;
        cursor: pointer;
        font-weight: bold;
    }
    .grid-button-click,
    .grid-button-red {
        background-color: gray;
        cursor: default;
    }
    .grid-button-red {
        background-color: red;
    }
    .grid-button:hover {
        background-color: silver;
    }
    .grid-button:active {
        background-color: darkgray;
    }
</style>
