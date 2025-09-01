fn main() -> Result<(), Box<dyn std::error::Error>> {
    unsafe { std::env::set_var("PROTOC", protoc_bin_vendored::protoc_bin_path()?) };

    let wkt_include = protoc_bin_vendored::include_path().unwrap();

    tonic_prost_build::configure().compile_protos(
        &["proto/events.proto"],
        &["proto", wkt_include.to_str().unwrap()],
    )?;

    Ok(())
}
