## ADDED Requirements

### Requirement: 通过扩展名检测视频URL
系统必须通过检查URL路径是否以.m3u8或.mp4结尾来识别视频URL。

#### Scenario: 检测到M3U8 URL
- **WHEN** 请求URL以.m3u8结尾
- **THEN** 系统标记该请求为视频内容

#### Scenario: 检测到MP4 URL
- **WHEN** 请求URL以.mp4结尾
- **THEN** 系统标记该请求为视频内容

#### Scenario: 忽略非视频URL
- **WHEN** 请求URL不以视频扩展名结尾
- **THEN** 系统不标记该请求为视频内容
