import {useState} from 'react';
import type {Snapshot} from './types';

// The underground market was a block inside the ledger.
//
// It does not belong there. The ledger answers "what am I worth and what is
// this costing me" — it is accounts. A price is not an account; it is a reason
// to go somewhere. Somebody checking whether moonshine is worth moving today is
// not doing bookkeeping, they are deciding where to spend the afternoon, and
// burying that under the daily upkeep made both jobs harder.
//
// So the market is its own page: what things cost against what they usually
// cost, what is in your hands, what holding it is costing you in attention, and
// where the trade is actually done.

const money = (n: number) => (n < 0 ? '−$' : '$') + Math.abs(Math.floor(n)).toLocaleString();

// Where each good changes hands. The core knows the trade rules; this is the
// one thing the player needs and cannot read off a price.
const counter: Record<string, string> = {
  moonshine: 'The waterfront deals in it. Pier 14 and the docks.',
  cigarettes: 'Mercer Exchange, openly enough.',
  arms: 'Nobody sells these across a counter. Ask at the garage.',
};

export function MarketScreen({world,onFind}: {world: Snapshot;onFind?:(id:string)=>void}) {
  const [apartmentFilter,setApartmentFilter]=useState('all'),[apartmentQuery,setApartmentQuery]=useState(''),[apartmentSort,setApartmentSort]=useState('address');
  const units=world.apartment_market||[],ownedUnits=units.filter(u=>u.owned);
  const apartmentListings=units.filter(u=>(apartmentFilter==='all'||apartmentFilter==='owned'&&u.owned||apartmentFilter==='available'&&u.available&&!u.owned)&&`${u.address} ${u.number} ${u.resident} ${u.owner}`.toLowerCase().includes(apartmentQuery.trim().toLowerCase())).sort((a,b)=>apartmentSort==='rent'?b.daily_rent-a.daily_rent||a.id.localeCompare(b.id):apartmentSort==='price'?(a.owned?a.offer:a.asking)-(b.owned?b.offer:b.asking)||a.id.localeCompare(b.id):a.address.localeCompare(b.address)||a.number-b.number);
  const goods = world.goods || [];
  const held = world.player.stock || {};
  const carrying = goods.filter(g => (held[g.id] || 0) > 0);
  // What the stock in your hands is costing you every day it stays there.
  const attention = carrying.reduce((sum, g) => sum + (held[g.id] || 0) * g.heat, 0);

  return (
    <section className="section-content">
      <div className="eyebrow">PRICES &amp; WHAT YOU ARE HOLDING</div>
      <h1 className="screen-title">Market notices</h1>
      {!!world.property_market?.length && <section className="property-exchange" aria-label="Property exchange">
        <h2>Property exchange</h2>
        <p>Deeds and standing broker offers. Visit an address to inspect the terms and complete a purchase or sale.</p>
        <div className="market-board">{world.property_market.map(property=><article className="market-good" key={property.id}>
          <header><b>{property.name}</b><span>{money(property.owned?property.offer:property.asking)}<small>{property.owned?' broker offer':property.available?' asking price':' reference price'}</small></span></header>
          <p>{property.owned?'Your deed':property.available?'Offered for sale':'Not currently offered'} · {property.condition}% condition</p>
          <p>{(property.neighborhood_index??100)<100 ? `Neighborhood prices ${100-(property.neighborhood_index??100)}% below normal after local violence. Quiet days help prices recover.` : 'Neighborhood prices are at their normal level.'}</p>
          <p>{property.residents} {property.home?'other ':''}{property.residents===1?'resident':'residents'}{property.home?' · Your current home':''}</p>
          <p>{property.owned?'A sale transfers the building and tenant accounts.':`Held by ${property.holder}.`}</p>
          {property.owned && property.home && <p>You can stay as a renter after selling; the terms are shown at the property.</p>}
          {onFind && <button className="plain" disabled={property.locked} onClick={()=>onFind(property.id)}>{property.locked?'District not yet accessible':'Inspect the address ↗'}</button>}
        </article>)}</div>
      </section>}
      {!!world.apartment_market?.length && <section className="property-exchange" aria-label="Apartment exchange">
        <h2>Apartment exchange</h2><p>Buy your home or a rental investment. Existing tenants stay; rent comes from their available cash. Vacant flats earn nothing until occupied. The broker lists a few homes at a time, with new listings as they sell.</p>
        <p className="apartment-holdings"><b>{ownedUnits.length} apartment {ownedUnits.length===1?'deed':'deeds'} held</b> · {money(ownedUnits.reduce((sum,u)=>sum+u.offer,0))} in current broker offers · {money(ownedUnits.reduce((sum,u)=>sum+(u.home?0:u.daily_rent),0))} scheduled rent per day. Actual collections depend on tenants’ cash.</p>
        <div className="apartment-filters">
          <label>Show<select value={apartmentFilter} onChange={e=>setApartmentFilter(e.target.value)}><option value="all">All listed apartments</option><option value="owned">Your deeds</option><option value="available">Available to buy</option></select></label>
          <label>Find an apartment<input type="search" value={apartmentQuery} onChange={e=>setApartmentQuery(e.target.value)} placeholder="Address, number, resident or owner"/></label>
          <label>Sort by<select value={apartmentSort} onChange={e=>setApartmentSort(e.target.value)}><option value="address">Address</option><option value="price">Price or offer · low first</option><option value="rent">Scheduled rent · high first</option></select></label>
        </div>
        <p role="status">{apartmentListings.length} of {units.length} listed apartments shown.</p>
        {apartmentListings.length===0&&<p>No apartments match these filters.</p>}
        <div className="market-board">{apartmentListings.map(unit=><article className="market-good" key={unit.id}>
          <header><b>{unit.address} · Apartment {unit.number}</b><span>{money(unit.owned?unit.offer:unit.asking)}<small>{unit.owned?' broker offer':unit.available?' asking price':' reference price'}</small></span></header>
          {unit.neighborhood_index!==undefined&&<p>{unit.neighborhood_index<100?`Neighborhood prices ${100-unit.neighborhood_index}% below normal after local violence. Quiet days help prices recover.`:'Neighborhood prices are at their normal level.'}</p>}
          <p>{unit.owned?'Your deed':`Owned by ${unit.owner}`} · {unit.home?'Your current home':unit.resident}</p>
          {(unit.owned || unit.available) && <p>{unit.home?(unit.owned?'No rent to pay':'Buying ends your rent'):unit.daily_rent>0?`${money(unit.daily_rent)} daily rent, collected from the resident’s available cash`:'No tenant income'}</p>}
          {onFind && <button className="plain" disabled={unit.locked} onClick={()=>onFind(unit.building)}>{unit.locked?'District not yet accessible':'Inspect the address ↗'}</button>}
        </article>)}</div>
      </section>}
      <h2>The underground market</h2>
      <p className="subtle market-note">
        Prices move whether or not anybody is watching them. Stock is only worth what somebody will
        pay for it today, and every day it sits in your hands it is drawing attention.
      </p>

      <div className="market-board">
        {goods.map(g => {
          const move = g.base ? Math.round(((g.price - g.base) / g.base) * 100) : 0;
          const have = held[g.id] || 0;
          return (
            <article key={g.id} className={'market-good' + (have ? ' holding' : '')}>
              <header>
                <b>{g.name}</b>
                <span className={move > 4 ? 'up' : move < -4 ? 'down' : ''}>
                  {money(g.price)}
                  <small> / {g.unit}</small>
                </span>
              </header>
              <p className="market-move">
                {move === 0
                  ? 'About what it usually goes for'
                  : `${Math.abs(move)}% ${move > 0 ? 'above' : 'below'} the usual ${money(g.base)}`}
              </p>
              <p className="market-holding">
                {have ? (
                  <>
                    <b>{have}</b> in your hands · {money(have * g.price)} at today's price
                  </>
                ) : (
                  <span className="subtle">none held</span>
                )}
              </p>
              <p className="market-where">
                {counter[g.id] || 'Traded quietly, where such things are traded.'}
              </p>
            </article>
          );
        })}
      </div>

      {carrying.length > 0 && (
        <div className="market-risk">
          <h2>What holding it costs</h2>
          <p>
            {carrying.map(g => `${held[g.id]} ${g.name.toLowerCase()}`).join(', ')} —{' '}
            <b>{attention}</b> attention a day, every day, until it moves. Stock found in a raid is
            stock lost, and the raid is likelier the longer it sits.
          </p>
        </div>
      )}

      {goods.length === 0 && (
        <p className="nothing-here">Nothing is being traded that you know of.</p>
      )}
    </section>
  );
}
