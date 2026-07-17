export namespace controller {
	
	export class ProgressInfo {
	    id: number;
	    name: string;
	    type: string;
	    done: number;
	    total: number;
	    percent: number;
	    timespec: number;
	
	    static createFrom(source: any = {}) {
	        return new ProgressInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.done = source["done"];
	        this.total = source["total"];
	        this.percent = source["percent"];
	        this.timespec = source["timespec"];
	    }
	}

}

export namespace model {
	
	export class Config {
	    max_concurrent_tasks: number;
	    max_segment_workers: number;
	    download_dir: string;
	    max_consecutive_errors: number;
	    default_headers: Record<string, string>;
	    interceptor_enabled: boolean;
	    agent_address: string;
	    vpn_address: string;
	    vpn_status: boolean;
	    gin_mode: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.max_concurrent_tasks = source["max_concurrent_tasks"];
	        this.max_segment_workers = source["max_segment_workers"];
	        this.download_dir = source["download_dir"];
	        this.max_consecutive_errors = source["max_consecutive_errors"];
	        this.default_headers = source["default_headers"];
	        this.interceptor_enabled = source["interceptor_enabled"];
	        this.agent_address = source["agent_address"];
	        this.vpn_address = source["vpn_address"];
	        this.vpn_status = source["vpn_status"];
	        this.gin_mode = source["gin_mode"];
	    }
	}
	export class Task {
	    id: number;
	    name: string;
	    url: string;
	    header: string;
	    type: string;
	    status: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.url = source["url"];
	        this.header = source["header"];
	        this.type = source["type"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

